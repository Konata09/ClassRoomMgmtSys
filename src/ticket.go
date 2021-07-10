package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type TicketOverview struct {
	Id             int    `json:"id,omitempty"`
	Title          string `json:"title,omitempty"`
	Severity       int    `json:"severity,omitempty"`
	Status         int    `json:"status,omitempty"`
	CreateUser     int    `json:"create_user,omitempty"`
	CreateUserName string `json:"create_user_name,omitempty"`
}

type Ticket struct {
	Id               int    `json:"id,omitempty"`
	Title            string `json:"title,omitempty"`
	Detail           string `json:"detail,omitempty"`
	Severity         int    `json:"severity,omitempty"`
	Status           int    `json:"status,omitempty"`
	ClassId          int    `json:"class_id,omitempty"`
	ClassroomName    string `json:"classroom_name,omitempty"`
	ClassroomGroup   string `json:"classroom_group,omitempty"`
	CreateUser       int    `json:"create_user,omitempty"`
	CreateUserName   string `json:"create_user_name,omitempty"`
	DutyUser1        int    `json:"duty_user_1,omitempty"`
	DutyUser1Name    string `json:"duty_user_1_name,omitempty"`
	DutyUser2        int    `json:"duty_user_2,omitempty"`
	DutyUser2Name    string `json:"duty_user_2_name,omitempty"`
	DutyUser3        int    `json:"duty_user_3,omitempty"`
	DutyUser3Name    string `json:"duty_user_3_name,omitempty"`
	CompleteUser     int    `json:"complete_user,omitempty"`
	CompleteUserName string `json:"complete_user_name,omitempty"`
	CreateTime       string `json:"create_time,omitempty"`
	StartTime        string `json:"start_time,omitempty"`
	CompleteTime     string `json:"complete_time,omitempty"`
	CompleteDetail   string `json:"complete_detail,omitempty"`
}

type AllTicket struct {
	Count   int              `json:"count,omitempty"`
	Tickets []TicketOverview `json:"tickets,omitempty"`
}

func AddTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var ticket Ticket
	json.NewDecoder(r.Body).Decode(&ticket)
	ret := addTicket(ticket.Title, ticket.Detail, ticket.Severity, ticket.ClassId, ticket.CreateUser, ticket.DutyUser1, ticket.DutyUser2, ticket.DutyUser3, ticket.CreateTime, ticket.StartTime)
	if ret {
		ApiOk(w)
		return
	} else {
		ApiErrMsg(w, "新建工单失败")
	}
}

func GetUserDutyTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	userId := r.URL.Query().Get("userid")
	if len(userId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(userId)
	if err != nil {
		ApiErrMsg(w, "用户id不正确")
		return
	}
	tickets := getUserDutyTicketOverview(id)
	json.NewEncoder(w).Encode(ApiReturn{0, "OK", AllTicket{
		Count:   len(tickets),
		Tickets: tickets,
	}})
}

func GetAllTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	countParam := r.URL.Query().Get("count")
	//fromParam := r.URL.Query().Get("from")
	//from, err := strconv.Atoi(fromParam)
	//if err != nil {
	//	from = 0
	//}
	count, err := strconv.Atoi(countParam)
	if err != nil {
		count = 20
	}
	if count < 1 {
		count = 20
	}
	tickets := getTickets(count)
	json.NewEncoder(w).Encode(ApiReturn{0, "OK", AllTicket{
		Count:   len(tickets),
		Tickets: tickets,
	}})

}
