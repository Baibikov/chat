package httperror

import (
	"encoding/json"
	"net/http"
)

type ErrMessage struct{
	Status 	int `json:"status"`
	Message string `json:"message"`
}

func New(w http.ResponseWriter, status int, message error) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ErrMessage{
		Message: message.Error(),
		Status: status,
	})
}
