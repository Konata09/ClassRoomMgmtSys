package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type ClassroomRedisStatus struct {
	ClassroomId   int            `json:"classroom_id"`
	ClassroomName string         `json:"classroom_name"` // 教室名称
	CourseName    string         `json:"course_name"`    // 课程名称
	TeacherName   string         `json:"teacher_name"`   // 教师姓名
	IsLive        int            `json:"is_live"`        // 是否直播
	IsRecord      int            `json:"is_record"`      // 是否录制及自动发布
	DeviceStatus  []DeviceStatus `json:"devices"`
}

type DeviceRedisStatus struct {
	DeviceId int `json:"device_id"`
	Ping     int `json:"ping"`   // ping ms
	Status   int `json:"status"` // 中控状态
}

func SetSingleClassroomStatusToRedis(status *ClassroomRedisStatus) {
	ctx := context.Background()
	marshal, _ := json.Marshal(status)

	_, err := rdb.Set(ctx, fmt.Sprintf("c%d", status.ClassroomId), marshal, 0).Result()
	if err != nil {
		logBoth("%s when SetSingleClassroomStatusToRedis ClassroomId: %d", err, status.ClassroomId)
	}
}

func GetSingleClassroomStatusFromRedis(classroomId int) *ClassroomRedisStatus {
	ctx := context.Background()
	res, err := rdb.Get(ctx, fmt.Sprintf("c%d", classroomId)).Result()
	if err != nil {
		logBoth("%s when GetSingleClassroomStatusFromRedis ClassroomId: %d", err, classroomId)
		return nil
	}
	var redisStatus ClassroomRedisStatus
	json.Unmarshal([]byte(res), &redisStatus)
	return &redisStatus
}

func SetMultiClassroomStatusToRedis(status []ClassroomRedisStatus) {
	ctx := context.Background()
	var classes []string
	for _, oneClass := range status {
		marshal, _ := json.Marshal(oneClass)
		classes = append(classes, fmt.Sprintf("c%d", oneClass.ClassroomId))
		classes = append(classes, string(marshal))
	}
	_, err := rdb.MSet(ctx, classes).Result()
	if err != nil {
		logBoth("%s when SetMultiClassroomStatusToRedis", err)
	}
}

func GetAllClassroomStatusFromRedis() []ClassroomRedisStatus {
	ctx := context.Background()
	classrooms := getClassrooms()
	var rooms []string
	for _, classroom := range classrooms {
		rooms = append(rooms, fmt.Sprintf("c%d", classroom.Id))
	}

	result, err := rdb.MGet(ctx, rooms...).Result()
	if err != nil {
		logBoth("%s when GetAllClassroomStatusFromRedis", err)
	}

	var redisStatusAll []ClassroomRedisStatus
	for i, classroom := range result {
		if classroom != nil {
			var redisStatus ClassroomRedisStatus
			err := json.Unmarshal([]byte(classroom.(string)), &redisStatus) // interface 转 string
			if err != nil {
				logBoth("%s when GetAllClassroomStatusFromRedis ClassroomId: %d", err, i+1)
				continue
			}
			redisStatusAll = append(redisStatusAll, redisStatus)
		} else {
			logBoth("GetAllClassroomStatusFromRedis return nil ClassroomId: %d", i+1)
		}
	}
	return redisStatusAll
}

func SetSingleDeviceStatusToRedis(status *DeviceRedisStatus) {
	ctx := context.Background()
	marshal, _ := json.Marshal(status)

	_, err := rdb.Set(ctx, fmt.Sprintf("d%d", status.DeviceId), marshal, 0).Result()
	if err != nil {
		logBoth("%s when SetSingleDeviceStatusToRedis DeviceId: %d", err, status.DeviceId)
	}
}

func GetSingleDeviceStatusFromRedis(deviceId int) *DeviceRedisStatus {
	ctx := context.Background()
	res, err := rdb.Get(ctx, fmt.Sprintf("d%d", deviceId)).Result()
	if err != nil {
		logBoth("%s when GetSingleDeviceStatusFromRedis DeviceId: %d", err, deviceId)
		return nil
	}
	var redisStatus DeviceRedisStatus
	json.Unmarshal([]byte(res), &redisStatus)
	return &redisStatus
}

func SetMultiDeviceStatusToRedis(status []DeviceRedisStatus) {
	ctx := context.Background()
	var devices []string
	for _, oneDevice := range status {
		marshal, _ := json.Marshal(oneDevice)
		devices = append(devices, fmt.Sprintf("d%d", oneDevice.DeviceId))
		devices = append(devices, string(marshal))
	}
	_, err := rdb.MSet(ctx, devices).Result()
	if err != nil {
		logBoth("%s when SetMultiDeviceStatusToRedis", err)
	}
}

func GetMultiDeviceStatusFromRedis(Ids []int) []DeviceRedisStatus {
	ctx := context.Background()
	var devs []string
	for _, id := range Ids {
		devs = append(devs, fmt.Sprintf("d%d", id))
	}

	result, err := rdb.MGet(ctx, devs...).Result()
	if err != nil {
		logBoth("%s when GetMultiDeviceStatusFromRedis", err)
	}

	var redisStatusAll []DeviceRedisStatus
	for i, device := range result {
		if device != nil {
			var redisStatus DeviceRedisStatus
			fmt.Printf("%v\n", device)
			err := json.Unmarshal([]byte(device.(string)), &redisStatus) // interface 转 string
			if err != nil {
				logBoth("%s when GetMultiDeviceStatusFromRedis DeviceId: %d", err, i+1)
				continue
			}
			redisStatusAll = append(redisStatusAll, redisStatus)
		} else {
			logBoth("GetMultiDeviceStatusFromRedis return nil DeviceId: %d", i+1)
		}
	}
	return redisStatusAll
}
