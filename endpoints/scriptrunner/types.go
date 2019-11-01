package scriptrunner

import ()

type script struct {
	Name           string
	Interpreter    string
	FilePath       string
	Stdin          string
	FileIn         string
	Description    string
	FileExtension  string
	UUID           string //this is not for the configuration. this is generated as it is sent to the client to keep track of which output stream goes to who.
	FileUpload     string
	CmdLineFlags   []flag
	CmdLineOptions []option
}

type flag struct {
	Name        string
	Prefix      string
	Description string
}

type option struct {
	Name        string
	Prefix      string
	Description string
}
