package scriptrunner

import (
	"encoding/json"
	"net/http"
	"strings"
)

//ScriptSearch is a naive search for a substring that returns bad results in a json.
func ScriptSearch(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	searchterm := r.FormValue("searchTerm")
	scripts := configuredScripts
	var ret []script

	for _, script := range scripts {
		//I don't know if two strings.contains is faster or if you concat them before the string.contains. I don't much care. In the future I'll probably do a indexof and send back information to highlight the matching substring so whats here right now will be in the trash soon enough.
		if strings.Contains(strings.ToLower(script.Name+script.Description), strings.ToLower(searchterm)) {
			ret = append(ret, script)
		}
	}
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(ret)
}
