package main

import (
	"encoding/json"
	"net/http"
)

type DutyCalender struct {
	Monday1        int    `json:"monday_1"`
	Monday2        int    `json:"monday_2"`
	Monday3        int    `json:"monday_3"`
	Monday1Name    string `json:"monday_1_name"`
	Monday2Name    string `json:"monday_2_name"`
	Monday3Name    string `json:"monday_3_name"`
	Tuesday1       int    `json:"tuesday_1"`
	Tuesday2       int    `json:"tuesday_2"`
	Tuesday3       int    `json:"tuesday_3"`
	Tuesday1Name   string `json:"tuesday_1_name"`
	Tuesday2Name   string `json:"tuesday_2_name"`
	Tuesday3Name   string `json:"tuesday_3_name"`
	Wednesday1     int    `json:"wednesday_1"`
	Wednesday2     int    `json:"wednesday_2"`
	Wednesday3     int    `json:"wednesday_3"`
	Wednesday1Name string `json:"wednesday_1_name"`
	Wednesday2Name string `json:"wednesday_2_name"`
	Wednesday3Name string `json:"wednesday_3_name"`
	Thursday1      int    `json:"thursday_1"`
	Thursday2      int    `json:"thursday_2"`
	Thursday3      int    `json:"thursday_3"`
	Thursday1Name  string `json:"thursday_1_name"`
	Thursday2Name  string `json:"thursday_2_name"`
	Thursday3Name  string `json:"thursday_3_name"`
	Friday1        int    `json:"friday_1"`
	Friday2        int    `json:"friday_2"`
	Friday3        int    `json:"friday_3"`
	Friday1Name    string `json:"friday_1_name"`
	Friday2Name    string `json:"friday_2_name"`
	Friday3Name    string `json:"friday_3_name"`
}

func SetDutyCalenderUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var dutyCal DutyCalender
	err := json.NewDecoder(r.Body).Decode(&dutyCal)
	if err != nil {
		ApiErr(w)
		return
	}
	setDutyCalender("Monday1", dutyCal.Monday1)
	setDutyCalender("Monday2", dutyCal.Monday2)
	setDutyCalender("Monday3", dutyCal.Monday3)
	setDutyCalender("Tuesday1", dutyCal.Tuesday1)
	setDutyCalender("Tuesday2", dutyCal.Tuesday2)
	setDutyCalender("Tuesday3", dutyCal.Tuesday3)
	setDutyCalender("Wednesday1", dutyCal.Wednesday1)
	setDutyCalender("Wednesday2", dutyCal.Wednesday2)
	setDutyCalender("Wednesday3", dutyCal.Wednesday3)
	setDutyCalender("Thursday1", dutyCal.Thursday1)
	setDutyCalender("Thursday2", dutyCal.Thursday2)
	setDutyCalender("Thursday3", dutyCal.Thursday3)
	setDutyCalender("Friday1", dutyCal.Friday1)
	setDutyCalender("Friday2", dutyCal.Friday2)
	setDutyCalender("Friday3", dutyCal.Friday3)
	ApiOk(w)
}

func GetDutyCalender(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//date := time.Now().Format("Monday")
	var dutyC *DutyCalender
	dutyC = getDutyCalender()
	json.NewEncoder(w).Encode(ApiReturn{0, "OK", dutyC})
}
