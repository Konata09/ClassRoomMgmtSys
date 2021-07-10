package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type SendUdpPacket struct {
	TargetDevId   []int  `json:"targetdevid"`
	Port          int    `json:"port"`
	CommandId     int    `json:"commandid"`
	UseCustom     bool   `json:"usecustom"`
	CustomPayload string `json:"custompayload"`
	Repeat        int    `json:"repeat"`
}

func sendSingleUdpPacket(ip string, port int, payload []byte) error {
	pc, err := net.ListenPacket("udp4", "")
	if err != nil {
		return errors.New(fmt.Sprintf("%s when sending packet to %s", err, ip))
	}

	defer pc.Close()
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return errors.New(fmt.Sprintf("%s when sending packet to %s", err, ip))
	}

	_, err = pc.WriteTo(payload, addr)
	if err != nil {
		return errors.New(fmt.Sprintf("%s when sending packet to %s", err, ip))
	}
	return nil
}

func SendUDP(body SendUdpPacket) ApiReturn {
	var payloadStr string
	var payloads [][]byte
	var command *Command
	var port int
	var repeat int
	if body.UseCustom {
		valid, msg := checkCommandValid("custom", body.CustomPayload, body.Port)
		if !valid {
			return ApiReturn{-1, msg, nil}
		}
		payloadStr = body.CustomPayload
		port = body.Port
	} else {
		command = getCommandById(body.CommandId)
		if command == nil {
			return ApiReturn{-1, "命令不存在", nil}
		}
		payloadStr = command.CommandValue
		if body.Port > 1 && body.Port < 65535 {
			port = body.Port
		} else {
			port = command.CommandPort
		}
	}
	for _, str := range strings.Split(payloadStr, ";") {
		hex, err := hexStringToByte(str)
		if err != nil {
			return ApiReturn{-1, err.Error(), nil}
		}
		payloads = append(payloads, hex)
	}
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
			dev := getDeviceById(tar)
			if dev == nil {
				appendMsg(&errMsg, fmt.Sprintf("id:%d 设备不存在", tar))
				continue
			}
			ip := dev.DeviceIp
			for _, hex := range payloads {
				err := sendSingleUdpPacket(ip, port, hex)
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
