package main

import (
	"fmt"
)

type SendWolPacket struct {
	TargetDevId       []int `json:"targetdevid"`
	DevAddress        bool  `json:"devaddress"`
	LocalNetBroadcast bool  `json:"localnetbroadcast"`
	SubNetBroadcast   bool  `json:"subnetbroadcast"`
	Repeat            int   `json:"repeat"`
}

var wolPort = 9

func SendWOL(body SendWolPacket) ApiReturn {
	if !body.DevAddress && !body.LocalNetBroadcast && !body.SubNetBroadcast {
		return ApiReturn{-1, "至少选择一个目的地址", nil}
	}

	var repeat int
	if body.Repeat == 0 {
		repeat = 1
	} else if body.Repeat > 5 {
		repeat = 5
	} else {
		repeat = body.Repeat
	}

	var errMsg string
	for i := 0; i < repeat; i++ {
		for _, tar := range body.TargetDevId {
			var address []string
			dev := getDeviceById(tar)
			if dev == nil {
				appendMsg(&errMsg, fmt.Sprintf("id:%d 设备不存在", tar))
				continue
			}
			if body.DevAddress {
				address = append(address, dev.DeviceIp)
			}
			if body.SubNetBroadcast {
				address = append(address, getSubnetBroadcast(dev.DeviceIp, 24))
			}
			if body.LocalNetBroadcast {
				address = append(address, "255.255.255.255")
			}
			payload, _ := hexStringToByte(getWolPayload(dev.DeviceMac))
			for _, addr := range address {
				err := sendSingleUdpPacket(addr, wolPort, payload)
				if err != nil {
					appendMsg(&errMsg, fmt.Sprintf("%s: %s", dev.DeviceName, err))
				}
			}
		}
	}
	if errMsg != "" {
		return ApiReturn{-1, errMsg, nil}
	}
	return ApiReturn{0, "OK", nil}
}
