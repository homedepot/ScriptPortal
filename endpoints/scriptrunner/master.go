package scriptrunner

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/homedepot/ScriptPortal/endpoints/httpHelpers"
	"github.com/homedepot/ScriptPortal/globals"
)

var singleScriptTemplate string = globals.TemplatePath + "inputform.html"

func startFlushingOutput(done chan bool, w http.Flusher) {
	for {
		select {
		case <-done:
			return
		default:
			time.Sleep(time.Second * 1)
			w.Flush()
		}
	}
}

func setNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Cache-Control", "post-check=0, pre-check=0")
}

func setupAndRunScript(w http.ResponseWriter, r *http.Request) {
	scriptList := configuredScripts
	scriptName := httpHelpers.GetRelativePath("script", r.URL.Path)
	targetScript := selectScript(scriptName, scriptList)
	setNoCacheHeaders(w)
	outputFileName := r.PostFormValue("outputFileName")
	//need a handling of unknown mime types more graceful than just not adding one
	mimeTypes, err := mime.ExtensionsByType(targetScript.FileExtension)
	if err == nil && len(mimeTypes) < 0 {
		w.Header().Set("Content-Type", mimeTypes[0])
	}
	contentDispositionString := fmt.Sprintf("attachment; filename=\"%s%s\"", outputFileName, targetScript.FileExtension)
	w.Header().Set("Content-Disposition", contentDispositionString)
	flusher, ok := w.(http.Flusher)
	flusher.Flush()
	if !ok {
		log.Fatal("oh crap")
	}

	var stdin bytes.Buffer

	done := make(chan bool)
	//this isn't strictly necessary, but it gets the download started quickly so the user realizes what it's doing. this can go away if there's an indication on the screen that the machine is working.
	go startFlushingOutput(done, flusher)
	defer func() { done <- true }()
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Println(err)
	}
	inputs := r.PostForm.Get("stdinLines")
	if inputs != "" {
		stdin.WriteString(inputs)
	}
	uuid := r.PostForm.Get("uuid")
	name, args := targetScript.buildScriptArguments(r.PostForm)

	if targetScript.FileIn == "true" {
		args = append(args, "-f")
		filenameForScript := getFileNameFromUpload(r.FormFile("UploadedFile"))

		//clean up after it's all been done
		defer os.Remove(filenameForScript)
		args = append(args, filenameForScript)
	}

	runScript(name, args, &stdin, w, uuid)
}

//ScriptMaster generates a scriptPage according to definitions in scriptconfig.json and runs the script given a POST request formatted properly
func ScriptMaster(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		scriptList := configuredScripts
		scriptName := httpHelpers.GetRelativePath("script", r.URL.Path)
		targetScript := selectScript(scriptName, scriptList)
		if targetScript.Name == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		targetScript.UUID = genUUID()
		templhtml, err := ioutil.ReadFile(globals.TemplatePath + "inputform.html")
		if err != nil {
			fmt.Print(err)
		}
		httpHelpers.AddContentToPage(w, string(templhtml), targetScript)
	case "POST":
		setupAndRunScript(w, r)

	default:
		w.WriteHeader(403)
		fmt.Fprintf(w, "This endpoint accepts only GET and POST options")
	}
}

func (target script) buildScriptArguments(form url.Values) (name string, arguments []string) {
	if target.Interpreter == "python" && runtime.GOOS == "linux" {
		target.Interpreter = "python3"
	}
	if target.Interpreter != "" {
		name = target.Interpreter
		arguments = append(arguments, target.FilePath)
	} else {
		name = target.FilePath
	}
	for _, option := range target.CmdLineOptions {
		if form.Get(option.Name) != "" {
			arguments = append(arguments, option.Prefix)
			arguments = append(arguments, form.Get(option.Name))
		}
	}
	for _, flag := range target.CmdLineFlags {
		if form.Get(flag.Name) == "true" {
			arguments = append(arguments, flag.Prefix)
		}

	}
	return
}

func runScript(scriptName string, args []string, input io.Reader, output io.Writer, uuid string) {
	cmd := exec.Command(scriptName, args...)
	cmd.Stdin = input
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}
	go io.Copy(output, stdout)

	go func(streamToUser io.Writer, commandOutput io.Reader) {
		scanner := bufio.NewScanner(commandOutput)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			streamToUser.Write(scanner.Bytes())
		}
	}(getOutputStream(uuid), stderr)

	err = cmd.Run()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}
