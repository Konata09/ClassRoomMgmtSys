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
	// 下面的字段只用于存储到 Redis 中
	CourseName  string `json:"course_name"`
	TeacherName string `json:"teacher_name"`
}

type AllClassroom struct {
	Count      int         `json:"count"`
	Classrooms []Classroom `json:"classrooms"`
}

type Camera struct {
	DeviceId  int    `json:"device_id"`
	Name      string `json:"name"`
	Type      int    `json:"device_type"`
	Ip        string `json:"ip"`
	Mac       string `json:"mac"`
	RtspAddr  string `json:"rtsp_addr"`
	RelayAddr string `json:"relay_addr"`
}

type ClassroomDetail struct {
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	GroupId   int      `json:"group_id"`
	GroupName string   `json:"group_name"`
	Devices   []Device `json:"devices"`
	Cameras   []Camera `json:"cameras"`
}

type ClassroomStatus struct {
	Id           int
	Rec          bool
	Live         bool
	CourseName   string
	TeacherName  string
	DeviceStatus []DeviceStatus
}

type DeviceStatus struct {
	Id     int
	Ping   int
	Status int
}

type SetClassroomJson struct {
	ClassId int    `json:"classid"`
	GroupId int    `json:"groupid"`
	Name    string `json:"name"`
}

func GetClassrooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Deprecated
	//var enPing = true
	//ParamPing := r.URL.Query().Get("ping")
	//if ParamPing == "false" {
	//	enPing = false
	//}

	classrooms := getClassrooms()
	redisStatus := GetAllClassroomStatusFromRedis()
	for i, classroom := range classrooms {
		for _, redisClass := range redisStatus {
			if classroom.Id == redisClass.ClassroomId {
				classrooms[i].Lindge = redisClass.Lindge
				classrooms[i].Controller = redisClass.Controller
				classrooms[i].Rec = redisClass.IsRecord != 0
				classrooms[i].Live = redisClass.IsLive != 0
				classrooms[i].TeacherName = redisClass.TeacherName
				classrooms[i].CourseName = redisClass.CourseName
				break
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
		if device.DeviceTypeId == 4 || device.DeviceTypeId == 5 || device.DeviceTypeId == 7 || device.DeviceTypeId == 8 || device.DeviceTypeId == 9 || device.DeviceTypeId == 10 {
			var cam Camera
			cam.DeviceId = device.DeviceId
			cam.Ip = device.DeviceIp
			cam.Mac = device.DeviceMac
			cam.Type = device.DeviceTypeId
			cam.Name = device.DeviceName
			switch device.DeviceTypeId {
			case 4: // Dahua
				cam.RtspAddr = fmt.Sprintf("rtsp://%s/cam/realmonitor?channel=1&subtype=0", device.DeviceIp)
				break
			case 7: // ZhiBo
				cam.RtspAddr = fmt.Sprintf("rtsp://%s", device.DeviceIp)
				break
			case 8: // ScreenEncoder
				cam.RtspAddr = fmt.Sprintf("rtsp://%s/stream", device.DeviceIp)
				break
			case 9: // Teacher Physical PTZ
				cam.RtspAddr = fmt.Sprintf("rtsp://%s/live/av0", device.DeviceIp)
				break
			case 10: // Student Physical PTZ
				cam.RtspAddr = fmt.Sprintf("rtsp://%s/live/av0", device.DeviceIp)
				break
			case 5: // Tiandy
				cam.RtspAddr = fmt.Sprintf("rtsp://%s/1/1", device.DeviceIp)
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
	go fetchSingleClassroomDeviceStatus(classId)
	redis := GetSingleClassroomStatusFromRedis(classId)
	var classroomStatus ClassroomStatus
	classroomStatus.DeviceStatus = redis.DeviceStatus
	classroomStatus.Live = redis.IsLive != 0
	classroomStatus.Rec = redis.IsRecord != 0
	classroomStatus.CourseName = redis.CourseName
	classroomStatus.TeacherName = redis.TeacherName
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
