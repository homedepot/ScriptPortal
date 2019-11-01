package scriptrunner

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

var configuredScripts []script

func Init(path string) (err error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return (err)
	}
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return (err)
	}
	json.Unmarshal(bytes, &configuredScripts)
	return
}

func selectScript(targetName string, list []script) (ret script) {
	for _, ret := range list {
		if ret.Name == targetName {
			return ret
		}
	}
	return
}

func genUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

//this is deceptively named. It actually does quite a bit more than just retrieving the file name.
func getFileNameFromUpload(file multipart.File, header *multipart.FileHeader, err error) (name string) {
	defer file.Close()
	fileNameParts := strings.Split(header.Filename, ".")

	name = "/tmp/scriptportal" + genUUID() + "." + fileNameParts[len(fileNameParts)-1]

	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	io.Copy(f, file)
	return
}
