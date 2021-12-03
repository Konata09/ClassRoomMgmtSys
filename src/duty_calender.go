package main

import (
	"encoding/json"
	"net/http"
)

type DutyUser struct {
	Id       string `json:"id"`
	Day      string `json:"day,omitempty"`
	Uid      int    `json:"uid,omitempty"`
	Username string `json:"username,omitempty"`
}

type DutyCalender struct {
	Monday    []DutyUser `json:"monday,omitempty"`
	Tuesday   []DutyUser `json:"tuesday,omitempty"`
	Wednesday []DutyUser `json:"wednesday,omitempty"`
	Thursday  []DutyUser `json:"thursday,omitempty"`
	Friday    []DutyUser `json:"friday,omitempty"`
	Saturday  []DutyUser `json:"saturday,omitempty"`
	Sunday    []DutyUser `json:"sunday,omitempty"`
}

func GetDutyCalender(w http.ResponseWriter, r *http.Request)  {
	var dutyC *DutyCalender
	dutyC = getDutyCalender()
	json.NewEncoder(w).Encode(ApiReturn{0, "OK", dutyC})
}

func SetDutyCalender(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var dutyUser DutyUser
		err := json.NewDecoder(r.Body).Decode(&dutyUser)
		if err != nil {
			ApiErr(w)
			return
		}
		if dutyUser.Day != "Monday" && dutyUser.Day != "Tuesday" && dutyUser.Day != "Wednesday" && dutyUser.Day != "Thursday" && dutyUser.Day != "Friday" && dutyUser.Day != "Saturday" && dutyUser.Day != "Sunday" {
			ApiErrMsg(w, "日期不存在")
			return
		}
		updateDutyCalender(dutyUser.Day, dutyUser.Uid, dutyUser.Id)
		ApiOk(w)
	case "PUT":
		var dutyUser DutyUser
		err := json.NewDecoder(r.Body).Decode(&dutyUser)
		if err != nil {
			ApiErr(w)
			return
		}
		if dutyUser.Day != "Monday" && dutyUser.Day != "Tuesday" && dutyUser.Day != "Wednesday" && dutyUser.Day != "Thursday" && dutyUser.Day != "Friday" && dutyUser.Day != "Saturday" && dutyUser.Day != "Sunday" {
			ApiErrMsg(w, "日期不存在")
			return
		}
		addDutyCalender(dutyUser.Day, dutyUser.Uid)
		ApiOk(w)
	case "DELETE":
		var dutyUser DutyUser
		err := json.NewDecoder(r.Body).Decode(&dutyUser)
		if err != nil {
			ApiErr(w)
			return
		}
		delDutyCalender(dutyUser.Id)
		ApiOk(w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}