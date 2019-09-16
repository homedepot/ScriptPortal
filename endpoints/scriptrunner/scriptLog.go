package scriptrunner

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/websocket"

	"github.com/homedepot/ScriptPortal/endpoints/httpHelpers"
)

var sockets map[string]myConn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type myConn struct {
	connection *websocket.Conn
}

func (c myConn) Write(message []byte) (n int, err error) {
	err = c.connection.WriteMessage(websocket.TextMessage, message)
	if err == nil {
		n = len(message)
	}
	return
}

//HandleSocketRequest handles socket requests..//todo: make a better description
func HandleSocketRequest(w http.ResponseWriter, r *http.Request) {
	//todo: make this next line something real
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	var socketContainer = myConn{ws}

	if err != nil {
		fmt.Println(err)
	}
	uuid := httpHelpers.GetRelativePath("status", r.URL.Path)
	if sockets == nil {
		sockets = map[string]myConn{uuid: socketContainer}
	} else {
		sockets[uuid] = socketContainer
	}
}

//todo: this seems a little too simple and easy. It should me more complicated
func getOutputStream(uuid string) (ret io.Writer) {
	stream, ok := sockets[uuid]
	if !ok {
		ret = os.Stdout
		return
	}
	stream.connection.WriteMessage(websocket.TextMessage, []byte("connected successfully to script logs"))
	ret = stream
	return
}
