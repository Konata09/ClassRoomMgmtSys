package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type SendUdpPacket struct {
	TargetId      []int  `json:"targetid"`
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

func SendUDP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var body SendUdpPacket
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			ApiErr(w)
			return
		}

		var payloadStr string
		var payloads [][]byte
		var command *Command
		var port int
		var repeat int
		if body.UseCustom {
			valid, msg := checkCommandValid("custom", body.CustomPayload, body.Port)
			if !valid {
				ApiErrMsg(w, msg)
				return
			}
			payloadStr = body.CustomPayload
			port = body.Port
		} else {
			command = getCommandById(body.CommandId)
			if command == nil {
				ApiErrMsg(w, "命令不存在")
				return
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
				ApiErrMsg(w, fmt.Sprint(err))
				return
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
			for _, tar := range body.TargetId {
				dev := getDeviceById(tar)
				if dev == nil {
					appendMsg(&errMsg, fmt.Sprintf("id:%d 设备不存在", tar))
					continue
				}
				if !dev.DeviceUdp {
					appendMsg(&errMsg, fmt.Sprintf("%s: 设备不支持发送普通UDP数据包", dev.DeviceName))
					continue
				}
				ip := dev.DeviceIp
				for _, hex := range payloads {
					err = sendSingleUdpPacket(ip, port, hex)
					if err != nil {
						appendMsg(&errMsg, fmt.Sprintf("%s: %s", dev.DeviceName, err))
					}
				}
			}
		}
		if errMsg != "" {
			ApiOkMsg(w, errMsg)
			return
		}
		ApiOk(w)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
