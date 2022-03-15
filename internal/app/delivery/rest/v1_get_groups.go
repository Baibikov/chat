package rest

import (
	"encoding/json"
	"net/http"

	"chat/internal/app/delivery/rest/types"
	"chat/pkg/httperror"
)

func (h *Handler) handlerV1GetGroups(w http.ResponseWriter, _ *http.Request) {
	groups := h.socket.GetGroups()
	resp := make(types.GroupsResponse, 0, len(groups))
	for _, g := range groups {
		resp = append(resp, types.Group{
			Name: g,
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	bb, err := json.Marshal(resp)
	if err != nil {
		httperror.New(w, http.StatusInternalServerError, err)
		return
	}

	_, err = w.Write(bb)
	if err != nil {
		httperror.New(w, http.StatusInternalServerError, err)
		return
	}
}