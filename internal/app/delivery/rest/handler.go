package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"chat/internal/app/delivery/service/socket"
)

type Handler struct {
	r *mux.Router

	upgrader websocket.Upgrader

	socket socket.Grouper
}

func New(socket socket.Grouper) *Handler {
	h := &Handler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  2<<9,
			WriteBufferSize: 2<<9,
		},
		socket: socket,
	}

	h.r = mux.NewRouter()

	return h
}

func (h *Handler) Setup() http.Handler {
	h.r.HandleFunc("/api/v1/chat/{groupName}", h.handlerWsV1GetMessage)
	h.r.HandleFunc("/api/v1/group", h.handlerV1GetGroups).Methods(http.MethodGet)
	h.r.HandleFunc("/api/v1/group", h.handlerV1CreateGroup).Methods(http.MethodPost)
	return h.r
}
