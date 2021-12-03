package main

import (
	"encoding/json"
	"net/http"
)

type DutyUser struct {
	Id       string `json:"id"`
	Day      string `json:"day"`
	Uid      int    `json:"uid"`
	Username string `json:"username"`
}

type DutyCalender struct {
	Monday    []DutyUser `json:"monday"`
	Tuesday   []DutyUser `json:"tuesday"`
	Wednesday []DutyUser `json:"wednesday"`
	Thursday  []DutyUser `json:"thursday"`
	Friday    []DutyUser `json:"friday"`
	Saturday  []DutyUser `json:"saturday"`
	Sunday    []DutyUser `json:"sunday"`
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