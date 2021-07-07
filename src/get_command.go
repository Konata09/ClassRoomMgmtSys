package main

import (
	"encoding/json"
	"net/http"
)

type UserCommand struct {
	CommandId   int    `json:"id"`
	CommandName string `json:"name"`
}
type CommandList struct {
	Count    int           `json:"count"`
	Commands []UserCommand `json:"commands"`
}

func GetCommand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		commands := getUserCommands()
		json.NewEncoder(w).Encode(&ApiReturn{
			Retcode: 0,
			Message: "OK",
			Data: &CommandList{
				Count:    len(commands),
				Commands: commands,
			},
		})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
