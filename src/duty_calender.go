package main

import (
	"encoding/json"
	"net/http"
)

//type DutyCalender struct {
//	Monday1        int    `json:"monday_1"`
//	Monday2        int    `json:"monday_2"`
//	Monday3        int    `json:"monday_3"`
//	Monday1Name    string `json:"monday_1_name"`
//	Monday2Name    string `json:"monday_2_name"`
//	Monday3Name    string `json:"monday_3_name"`
//	Tuesday1       int    `json:"tuesday_1"`
//	Tuesday2       int    `json:"tuesday_2"`
//	Tuesday3       int    `json:"tuesday_3"`
//	Tuesday1Name   string `json:"tuesday_1_name"`
//	Tuesday2Name   string `json:"tuesday_2_name"`
//	Tuesday3Name   string `json:"tuesday_3_name"`
//	Wednesday1     int    `json:"wednesday_1"`
//	Wednesday2     int    `json:"wednesday_2"`
//	Wednesday3     int    `json:"wednesday_3"`
//	Wednesday1Name string `json:"wednesday_1_name"`
//	Wednesday2Name string `json:"wednesday_2_name"`
//	Wednesday3Name string `json:"wednesday_3_name"`
//	Thursday1      int    `json:"thursday_1"`
//	Thursday2      int    `json:"thursday_2"`
//	Thursday3      int    `json:"thursday_3"`
//	Thursday1Name  string `json:"thursday_1_name"`
//	Thursday2Name  string `json:"thursday_2_name"`
//	Thursday3Name  string `json:"thursday_3_name"`
//	Friday1        int    `json:"friday_1"`
//	Friday2        int    `json:"friday_2"`
//	Friday3        int    `json:"friday_3"`
//	Friday1Name    string `json:"friday_1_name"`
//	Friday2Name    string `json:"friday_2_name"`
//	Friday3Name    string `json:"friday_3_name"`
//}

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