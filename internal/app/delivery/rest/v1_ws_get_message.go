package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"chat/internal/app/delivery/rest/types"
	"chat/internal/app/delivery/service/socket"
	"chat/pkg/httperror"
)


func (h *Handler) handlerWsV1GetMessage(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["groupName"]

	h.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		httperror.New(w, http.StatusInternalServerError, err)
		return
	}

	conn, err := h.socket.CreateConnect(socket.GroupName(groupName), ws)
	if err != nil {
		httperror.New(w, http.StatusInternalServerError, err)
		return
	}

	err = conn.Read(r.Context(), func(p *socket.Payload) error {
		var rr types.ChatMessageRequest
		err := json.Unmarshal(p.Body, &rr)
		if err != nil {
			return err
		}

		p.Write(rr)
		return nil
	})
	if err != nil {
		httperror.New(w, http.StatusInternalServerError, err)
		return
	}
}