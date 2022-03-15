package types

import "encoding/json"

type GroupRequest struct {
	Name string `json:"name"`
}

type ChatMessageResponse struct {
	UserID  json.Number `json:"userId"`
	Message string `json:"message"`
}

type ChatMessageRequest struct {
	UserID json.Number 	`json:"userId"`
	Message string 		`json:"message"`
}

type GroupsResponse []Group

type Group struct {
	Name string `json:"name"`
}