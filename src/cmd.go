package main

import (
	"encoding/json"
	"net/http"
)

/*
CmdId - CmdName
	1 - 中控开
	2 - 中控关
	3 - 云盒开 WOL
	4 - 云盒关
	5 - 投影云盒
	6 - 幕布升
	7 - 幕布停
	8 - 幕布降
	9 - 音量加
	10 - 音量减
	11 - 静音
	12 - 静音取消
	13 - 投影开
	14 - 投影关
*/
type SendCmd struct {
	ClassId int
	CmdName string
	CmdId   int
}

func HandleCmd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var sendCmdJson SendCmd
	json.NewDecoder(r.Body).Decode(&sendCmdJson)
	classroom := getClassroom(sendCmdJson.ClassId)
	var targetId int
	var ret ApiReturn
	if sendCmdJson.CmdId == 3 {
		for _, dev := range classroom.Devices {
			if dev.DeviceTypeId == 2 {
				targetId = dev.DeviceId
				break
			}
		}
		ret = SendWOL(SendWolPacket{
			TargetDevId:       []int{targetId},
			DevAddress:        true,
			LocalNetBroadcast: true,
			SubNetBroadcast:   true,
			Repeat:            3,
		})
	} else {
		for _, dev := range classroom.Devices {
			if dev.DeviceTypeId == 1 {
				targetId = dev.DeviceId
				break
			}
		}
		ret = SendUDP(SendUdpPacket{
			TargetDevId:   []int{targetId},
			Port:          0,
			CommandId:     sendCmdJson.CmdId,
			UseCustom:     false,
			CustomPayload: "",
			Repeat:        1,
		})
	}
	json.NewEncoder(w).Encode(&ret)
}
