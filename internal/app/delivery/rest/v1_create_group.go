package rest

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"chat/internal/app/delivery/rest/types"
	"chat/internal/app/delivery/service/socket"
	"chat/pkg/httperror"
)

func (h *Handler) handlerV1CreateGroup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var group types.GroupRequest
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		httperror.New(w, http.StatusBadRequest, err)
		return
	}

	err = h.socket.CreateGroup(socket.GroupName(group.Name))
	if err != nil {
		httperror.New(w, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("make group success %s", group.Name)
}
