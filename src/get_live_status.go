package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type Reserve struct {
	RoomName      string `json:"room-name,omitempty"`       // 教室名称
	Name          string `json:"name,omitempty"`            // 课程名称
	TeacherName   string `json:"teacher-name,omitempty"`    // 教师姓名
	ReserveStatus int    `json:"reserve-status,omitempty"`  // 1:进行中 0:未开始 2:已完成
	IsLive        int    `json:"is-live,omitempty"`         // 是否直播
	IsRecordFile  int    `json:"is-record-file,omitempty"`  // 是否录制
	IsAutoPublish int    `json:"is-auto-publish,omitempty"` // 是否自动发布
}

var sessionID = ""

func dyLogin() bool {
	url := "http://172.31.163.248:8885/GatewayCenter/users/login"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("user-id", "admin")
	_ = writer.WriteField("password", "Jnxq123456")
	_ = writer.WriteField("remember", "0")
	_ = writer.WriteField("tenant-id", "default")
	_ = writer.WriteField("page-from", "1")
	writer.Close()

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, payload)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		logBoth("[ERR] %s when dyLogin DoRequest", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logBoth("[ERR] %s when dyLogin ReadAll", err)
	}
	if strings.Contains(string(body), "系统管理员") {
		cookies := res.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "SESSION" {
				sessionID = cookie.Value
				break
			}
		}
		logBoth("[INFO] dyLogin Success. SessionID is %s", sessionID)
		return true
	} else {
		logBoth("[ERR] dyLogin Failed. %s", string(body))
		return false
	}
}

func getLiveStatusFromDy() []Reserve {
	url := "http://172.31.163.248:8885/GatewayCenter/reserves?reserve-status=1&count=50&start=0"
	method := "GET"

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)

	req.Header.Add("Cookie", sessionID)
	req.Header.Add("Cache-Control", "no-cache")

	client.Timeout = time.Second * 2
	res, err := client.Do(req)
	if err != nil {
		logBoth("[ERR] %s when getLiveStatusFromDy DoRequest", err)
		return nil
	}

	defer res.Body.Close()
	if err != nil {
		logBoth("[ERR] %s when getLiveStatusFromDy ReadAll", err)
		return nil
	}
	var reserves []Reserve
	err = json.NewDecoder(res.Body).Decode(&reserves)
	if err != nil {
		logBoth("[ERR] %s when getLiveStatusFromDy DecodeJson %s", err, res.Body)
		return nil
	}
	fmt.Println(reserves)
	return reserves
}
