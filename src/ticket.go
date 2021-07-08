package main

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
