package httpHelpers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/homedepot/ScriptPortal/globals"
)

const boilerPlateFileName string = "BoilerPlate.html.template"

var Navs []TopBarInfo

type TopBarInfo struct {
	Name, Link string
}

type contentWrapper struct {
	Data        string
	Navigations []TopBarInfo
}

//this and the next function are bad and I need to re-do how the templating system works in the app.
func WrapInBoilerPlate(w http.ResponseWriter, data string) {
	temp := template.Must(template.ParseFiles(globals.TemplatePath + boilerPlateFileName))
	err := temp.Execute(w, contentWrapper{data, Navs})
	if err != nil {
		fmt.Printf("this %s", err)
	}
}

//AddContentToPage just opens a template, sends whatever data it's given to it, and adds that to a http response writer
func AddContentToPage(w http.ResponseWriter, contentTemplate string, data interface{}) {
	var buf bytes.Buffer

	template.Must(template.New("").Parse(contentTemplate)).Execute(&buf, data)
	things := buf.Bytes()
	WrapInBoilerPlate(w, string(things))
}

//GetRelativePath takes a path and gives you the path after the root. there's more algorithm to it than I thought, so i made a helper func
func GetRelativePath(root string, path string) (ret string) {
	ret = strings.TrimSuffix(path, ".html")
	ret = strings.Trim(ret, "/")
	if len(ret) <= len(root) {
		ret = ""
		return
	}
	ret = ret[len(root):]
	ret = strings.Trim(ret, "/")
	return
}
