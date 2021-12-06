package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type ClassroomRedisStatus struct {
	ClassroomId   int            `json:"classroom_id"`
	ClassroomName string         `json:"classroom_name"`  // 教室名称
	CourseName    string         `json:"course_name"`     // 课程名称
	TeacherName   string         `json:"teacher_name"`    // 教师姓名
	ReserveStatus int            `json:"reserve_status"`  // 1:进行中 0:未开始 2:已完成
	IsLive        int            `json:"is_live"`         // 是否直播
	IsRecordFile  int            `json:"is_record_file"`  // 是否录制
	IsAutoPublish int            `json:"is_auto_publish"` // 是否自动发布
	DeviceStatus  []DeviceStatus `json:"devices"`
}

func SetSingleStatusToRedis(status *ClassroomRedisStatus) {
	ctx := context.Background()
	marshal, _ := json.Marshal(status)

	_, err := rdb.Set(ctx, fmt.Sprintf("c%d", status.ClassroomId), marshal, 0).Result()
	if err != nil {
		logBoth("%s when SetSingleStatusToRedis ClassroomId: %d", err, status.ClassroomId)
	}
}

func GetSingleStatusFromRedis(classroomId int) *ClassroomRedisStatus {
	ctx := context.Background()
	res, err := rdb.Get(ctx, fmt.Sprintf("c%d", classroomId)).Result()
	if err != nil {
		logBoth("%s when GetSingleStatusFromRedis ClassroomId: %d", err, classroomId)
		return nil
	}
	var redisStatus ClassroomRedisStatus
	json.Unmarshal([]byte(res), &redisStatus)
	return &redisStatus
}
func SetAllStatusFromRedis(status []ClassroomRedisStatus) {
	ctx := context.Background()
	var classes []string
	for _, oneClass := range status {
		marshal, _ := json.Marshal(oneClass)
		classes = append(classes, fmt.Sprintf("c%d", oneClass.ClassroomId))
		classes = append(classes, string(marshal))
	}
	_, err := rdb.MSet(ctx, classes).Result()
	if err != nil {
		logBoth("%s when SetAllStatusFromRedis", err)
	}
}

func GetAllStatusFromRedis() []ClassroomRedisStatus {
	ctx := context.Background()
	classrooms := getClassrooms()
	var rooms []string
	for _, classroom := range classrooms {
		rooms = append(rooms, fmt.Sprintf("c%d", classroom.Id))
	}

	result, err := rdb.MGet(ctx, rooms...).Result()
	if err != nil {
		logBoth("%s when GetAllStatusFromRedis", err)
	}

	var redisStatusAll []ClassroomRedisStatus
	for i, classroom := range result {
		var redisStatus ClassroomRedisStatus
		if classroom != nil {
			fmt.Printf("%v\n", classroom)
			err := json.Unmarshal([]byte(classroom.(string)), &redisStatus) // interface 转 string
			if err != nil {
				logBoth("%s when GetAllStatusFromRedis ClassroomId: %d", err, i+1)
				continue
			}
			redisStatusAll = append(redisStatusAll, redisStatus)
		} else {
			logBoth("GetAllStatusFromRedis return nil ClassroomId: %d", i+1)
		}
	}
	return redisStatusAll
}
