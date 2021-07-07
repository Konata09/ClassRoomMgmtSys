package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

type Command struct {
	CommandId    int    `json:"id"`
	CommandName  string `json:"name"`
	CommandValue string `json:"value"`
	CommandPort  int    `json:"port"`
}
type AllCommands struct {
	Count    int       `json:"count"`
	Commands []Command `json:"commands"`
}

func checkCommandValid(name string, value string, port int) (bool, string) {
	reCmd := regexp.MustCompile(`^[A-Fa-f0-9;]+$|^$`)

	if strings.TrimSpace(name) == "" || strings.TrimSpace(value) == "" {
		return false, "命令名称或者内容为空"
	}
	if !reCmd.MatchString(value) {
		return false, "命令格式不正确"
	}
	if port < 1 || port > 65535 {
		return false, "端口需位于1-65535之间"
	}
	return true, ""
}

func SetCommand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		commands := getCommands()
		json.NewEncoder(w).Encode(&ApiReturn{
			Retcode: 0,
			Message: "OK",
			Data: &AllCommands{
				Count:    len(commands),
				Commands: commands,
			},
		})
	case "PUT":
		var body AllCommands
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			ApiErr(w)
			return
		}
		var commands []Command
		var msg string
		for _, cmd := range body.Commands {
			valid, m := checkCommandValid(cmd.CommandName, cmd.CommandValue, cmd.CommandPort)
			if valid {
				commands = append(commands, cmd)
			} else {
				msg = msg + m + " "
			}
		}
		if len(commands) == 0 {
			ApiErrMsg(w, msg+"No item to add")
			return
		}
		ok := addCommand(commands)
		if ok {
			ApiOkMsg(w, msg+"OK")
		} else {
			ApiErrMsg(w, msg+"请求错误")
		}
	case "POST":
		var body Command
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			ApiErr(w)
			return
		}
		if getCommandById(body.CommandId) == nil {
			ApiErrMsg(w, "命令不存在")
			return
		}
		valid, m := checkCommandValid(body.CommandName, body.CommandValue, body.CommandPort)
		if !valid {
			ApiErrMsg(w, m)
			return
		}
		ok := setCommand(body.CommandId, body.CommandName, body.CommandValue, body.CommandPort)
		if ok {
			ApiOk(w)
		} else {
			ApiErr(w)
		}
	case "DELETE":
		var body Command
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			ApiErr(w)
			return
		}
		if getCommandById(body.CommandId) == nil {
			ApiErrMsg(w, "命令不存在")
			return
		}
		ok := deleteCommand(body.CommandId)
		if ok {
			ApiOk(w)
		} else {
			ApiErr(w)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
