package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
	Upgrader Upgrader
}

func (handler *Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	conn, err := handler.Upgrader.Upgrade(response, request, response.Header())
	if err != nil {
		// TODO: Logging?
		return
	}
	conn.Close()
}

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}
