package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Classroom struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	GroupId    int    `json:"group_id"`
	GroupName  string `json:"group_name"`
	Controller int    `json:"controller"`
	Lindge     int    `json:"lindge"`
	Live       bool   `json:"live"`
	REC        bool   `json:"rec"`
}

type AllClassroom struct {
	Count      int         `json:"count,omitempty"`
	Classrooms []Classroom `json:"classrooms,omitempty"`
}

func GetClassrooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	classrooms := getClassrooms()
	devices := getClassroomsMainDevices()
	pingDevices(devices)
	fmt.Println(devices)
	json.NewEncoder(w).Encode(&ApiReturn{
		Retcode: 0,
		Message: "OK",
		Data: AllClassroom{
			Count:      len(classrooms),
			Classrooms: classrooms,
		},
	})
}
