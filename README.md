# ScriptPortal
![](https://github.com/homedepot/scriptportal/workflows/Go/badge.svg)

## What is it?
ScriptPortal is a webapp that wraps the exec function in golang to allow users to run in house created scripts at their own leisure

### What problem does it solve?
Utilities that were previously reserved for the command line officionados are now effectively tiny web apps integrated into a neat little portal. Additionally anything that would be slightly too small to be a full-fledged web app can be integrated via the plugin system with little effort, and smoothly appear side-by-side with the main system.  

In addition, the data flow of the application allows for users to see what's happening in the program in real time while the job is being run as opposed to a schedule and pick up later workflow that would make it difficult to see what is happening or if there were any issues in the execution of the target program. 

## Getting started
1. clone the repo (preferably using go get)
2. cd to the repo
3. run the docker build command ```docker build -t scriptportal .```
4. run the container and bind port 80 on local to the container ```docker run -p80:80 scriptportal```
5. connect to [localhost](http://localhost) from your web browser

## Script Configuration:
The format is defined by the following go code, encoded according to encoding/json standards. 
```go
type script struct {
        Name           string // The name of the script to run
        Interpreter    string // The interpreter set to run the script (optional if calling a binary file)
        FilePath       string // The full path of the target file to be run
        Stdin          string // true if input is taken on standard input. This will be accessed by the user as a text input box on the webpage for that script
        FileIn         string // true if there is a file to be input. The script must accept a -f flag followed by the path of the file to be input for FileIn to work properly. The name of the file uploaded will be changed to a randomly generated one and passed via -f argument to your script
        Description    string // to be shown to the user. Usually gives a short explanation of the inputs and outputs and why you may wish to run the script
        FileExtension  string // the file extension of the output file. There to set the mime type of the response. 
        FileUpload     string // present if there is a file 
	OutputType     string // either "stdout" or "file" depending on if the program writes to standard output(prefered) or in a file format that cannot be expressed via stdout
        CmdLineFlags   []flag //defined below
        CmdLineOptions []option //defined below
}
//definition of command line flags given to the program. name and description are shown to the user as well as a checkbox to indicate wether the option is there or not. prefix is the flag (including the dash) that is handed to the program as an argument
type flag struct {
        Name        string
        Prefix      string
        Description string
}
//Same as flag, but shown to the user as a text box. The text the user puts in is the next argument after the prefix argument when passed to the target script.
type option struct {
        Name        string
        Prefix      string
        Description string
}

```

## Plugins:
The plugin functionality is for any additional functionality you wish to add to the portal that may otherwise be too small for it's own standalone web app.
plugin spec:

Each plugin must be go file placed in the scriptPortal/plugins/ directory. The build script compiles each .go file in the directory to a .so file, which is moved into the /var/include/ScriptPortal when the install script is run. At runtime, the main program will link the so files and add endpoints for each plugin. Each plugin must be written with the following exported functions

Name -> a function that returns type string of the name of the plugin at the top of the page

ex: 
```go
func Name() string { return "abcfunctionality" }
```
Path -> a function that returns type string of the link to be matched by the http handlefunc

ex: 
```go	
func Path() string { return "/abc/" }
```
	
HandleRequest -> a function that takes types http.ResponseWriter, *http.Request, and chan []byte. If the endpoint is meant to have information that shows within the header and footer information, it should send type []byte of the data to be included over the channel, and then close it. If the data is to be added directly to the http.ResponseWriter, write to it like a normal http response and close the channel when you are done.

ex: 
```go
func HandleRequest(w http.ResponseWriter, r *http.Request, resChan chan []byte) {
	defer close(resChan)
	if (some condition){
		w.write(an_api_type_response)
	}else{
		resChan <- []byte("Sorry pal, I don't like your tone")
}
```
