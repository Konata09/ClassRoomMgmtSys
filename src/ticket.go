package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type TicketOverview struct {
	Id             int    `json:"id"`
	Title          string `json:"title"`
	Severity       int    `json:"severity"`
	Status         int    `json:"status"`
	CreateUser     int    `json:"create_user"`
	CreateUserName string `json:"create_user_name"`
}

type Ticket struct {
	Id               int    `json:"id"`
	Title            string `json:"title"`
	Detail           string `json:"detail"`
	Severity         int    `json:"severity"`
	Status           int    `json:"status"`
	ClassId          int    `json:"class_id"`
	ClassroomName    string `json:"classroom_name"`
	ClassroomGroup   string `json:"classroom_group"`
	CreateUser       int    `json:"create_user"`
	CreateUserName   string `json:"create_user_name"`
	DutyUser1        int    `json:"duty_user_1"`
	DutyUser1Name    string `json:"duty_user_1_name"`
	DutyUser2        int    `json:"duty_user_2"`
	DutyUser2Name    string `json:"duty_user_2_name"`
	DutyUser3        int    `json:"duty_user_3"`
	DutyUser3Name    string `json:"duty_user_3_name"`
	CompleteUser     int    `json:"complete_user"`
	CompleteUserName string `json:"complete_user_name"`
	CreateTime       string `json:"create_time"`
	StartTime        string `json:"start_time"`
	CompleteTime     string `json:"complete_time"`
	CompleteDetail   string `json:"complete_detail"`
}

type AllTicket struct {
	Count   int              `json:"count"`
	Tickets []TicketOverview `json:"tickets"`
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
func GetTicketDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		ApiErrMsg(w, "id 无效")
		return
	}
	ticket := getTicket(id)
	if ticket == nil {
		ApiErr(w)
		return
	}
	json.NewEncoder(w).Encode(ApiReturn{0, "OK", ticket})
}
