package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Classroom struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	GroupId    int    `json:"group_id"`
	GroupName  string `json:"group_name"`
	Controller int    `json:"controller"`
	Lindge     int    `json:"lindge"`
	Live       bool   `json:"live"`
	Rec        bool   `json:"rec"`
}

type AllClassroom struct {
	Count      int         `json:"count,omitempty"`
	Classrooms []Classroom `json:"classrooms,omitempty"`
}
type Camera struct {
	DeviceId  int    `json:"device_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      int    `json:"device_type,omitempty"`
	Ip        string `json:"ip,omitempty"`
	Mac       string `json:"mac,omitempty"`
	RtspAddr  string `json:"rtsp_addr,omitempty"`
	RelayAddr string `json:"relay_addr,omitempty"`
}

type ClassroomDetail struct {
	Id        int      `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	GroupId   int      `json:"group_id,omitempty"`
	GroupName string   `json:"group_name,omitempty"`
	Devices   []Device `json:"devices,omitempty"`
	Cameras   []Camera `json:"cameras,omitempty"`
}

type ClassroomStatus struct {
	Id           int
	Rec          bool
	Live         bool
	DeviceStatus []DeviceStatus
}

type DeviceStatus struct {
	Id     int
	Ping   int
	Status int
}

type SetClassroomJson struct {
	ClassId int    `json:"classid,omitempty"`
	GroupId int    `json:"groupid,omitempty"`
	Name    string `json:"name,omitempty"`
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

func GetClassroomDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	ParamClassid := r.URL.Query().Get("classid")
	var classId int
	if len(ParamClassid) > 0 {
		var err error
		classId, err = strconv.Atoi(ParamClassid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	classroomDetail := getClassroom(classId)
	var cameras []Camera
	for _, device := range classroomDetail.Devices {
		if device.DeviceTypeId == 4 || device.DeviceTypeId == 5 || device.DeviceTypeId == 7 || device.DeviceTypeId == 8 {
			var cam Camera
			cam.DeviceId = device.DeviceId
			cam.Ip = device.DeviceIp
			cam.Mac = device.DeviceMac
			cam.Type = device.DeviceTypeId
			cam.Name = device.DeviceName
			switch device.DeviceTypeId {
			case 4: // Dahua
				cam.RtspAddr = fmt.Sprintf("rtsp://%s", device.DeviceIp)
				break
			case 5: // Tiandy
				cam.RtspAddr = fmt.Sprintf("rtsp://%s", device.DeviceIp)
				break
			case 7: // ZhiBo
				cam.RtspAddr = fmt.Sprintf("rtsp://%s", device.DeviceIp)
				break
			case 8: // ScreenEncoder
				cam.RtspAddr = fmt.Sprintf("rtsp://%s", device.DeviceIp)
				break
			}
			cameras = append(cameras, cam)
		}
	}
	classroomDetail.Cameras = cameras
	json.NewEncoder(w).Encode(&ApiReturn{
		Retcode: 0,
		Message: "OK",
		Data:    classroomDetail,
	})
}

func GetClassroomStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	ParamClassid := r.URL.Query().Get("classid")
	var classId int
	if len(ParamClassid) > 0 {
		var err error
		classId, err = strconv.Atoi(ParamClassid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	doneDevice := make(chan int)
	doneController := make(chan DetectRes)
	devices := getDevicesByClassId(classId)
	var controllerIp string
	for _, dev := range devices {
		if dev.DeviceTypeId == 1 {
			controllerIp = dev.DeviceIp
			break
		}
	}
	go pingDevices(devices, doneDevice)
	go getControollerStatusSingle(controllerIp, 1, doneController)
	var classroomStatus ClassroomStatus
	classroomStatus.Id = classId
	var devStatus []DeviceStatus
	controllerRes := <-doneController
	if <-doneDevice == 1 {
		for _, dev := range devices {
			var devStat DeviceStatus
			devStat.Status = dev.status
			devStat.Ping = dev.pingRes
			devStat.Id = dev.DeviceId
			if dev.DeviceTypeId == 1 {
				devStat.Status = controllerRes.res
			}
			devStatus = append(devStatus, devStat)
		}
	}
	classroomStatus.DeviceStatus = devStatus
	json.NewEncoder(w).Encode(&ApiReturn{
		Retcode: 0,
		Message: "OK",
		Data:    classroomStatus,
	})
}

func SetClassroom(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	var setJson SetClassroomJson
	err := json.NewDecoder(r.Body).Decode(&setJson)
	if err != nil {
		ApiErr(w)
		return
	}
	classid := setJson.ClassId
	groupid := setJson.GroupId
	classname := setJson.Name
	ret := setClassroom(classid, classname, groupid)
	if ret {
		ApiOk(w)
	} else {
		ApiErrMsg(w, "修改失败")
	}
}
