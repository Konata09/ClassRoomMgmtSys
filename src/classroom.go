package main

import (
	"encoding/json"
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
	var enPing = true
	ParamPing := r.URL.Query().Get("ping")
	if ParamPing == "false" {
		enPing = false
	}
	classrooms := getClassrooms()
	if enPing {
		done := make(chan int)
		lindges := getClassroomLindges()
		controllers := getClassroomControllers()
		go pingDevices(lindges, done)              // 1
		go getControllersStatus(controllers, done) // 2
		doneOk1 := <-done
		if doneOk1 == 1 {
			for _, lindge := range lindges {
				for i, classroom := range classrooms {
					if classroom.Id == lindge.DeviceClassId {
						classrooms[i].Lindge = lindge.pingRes
						break
					}
				}
			}
		} else if doneOk1 == 2 {
			for _, controller := range controllers {
				for i, classroom := range classrooms {
					if classroom.Id == controller.DeviceClassId {
						classrooms[i].Controller = controller.status
						break
					}
				}
			}
		}
		doneOk2 := <-done
		if doneOk2 == 1 {
			for _, lindge := range lindges {
				for i, classroom := range classrooms {
					if classroom.Id == lindge.DeviceClassId {
						classrooms[i].Lindge = lindge.pingRes
						break
					}
				}
			}
		} else if doneOk2 == 2 {
			for _, controller := range controllers {
				for i, classroom := range classrooms {
					if classroom.Id == controller.DeviceClassId {
						classrooms[i].Controller = controller.status
						break
					}
				}
			}
		}
	}
	json.NewEncoder(w).Encode(&ApiReturn{
		Retcode: 0,
		Message: "OK",
		Data: AllClassroom{
			Count:      len(classrooms),
			Classrooms: classrooms,
		},
	})
}
