package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

type Device struct {
	DeviceId      int    `json:"id"`
	DeviceName    string `json:"name"`
	DeviceIp      string `json:"ip"`
	DeviceMac     string `json:"mac"`
	DeviceUdp     bool   `json:"udp"`
	DeviceWol     bool   `json:"wol"`
	DeviceSubmask int    `json:"submask"`
}

type AllDevices struct {
	Count   int      `json:"count"`
	Devices []Device `json:"devices"`
}

func checkDeviceValid(name string, ip string, mac string, udp bool, wol bool, submask int) (bool, string) {
	reIp := regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4}$`)
	reMac := regexp.MustCompile(`^([0-9A-Fa-f]{2}([:-]?|$)){6}$`)

	if !udp && !wol {
		return false, "WOL与UDP中必须选择一项"
	}
	if strings.TrimSpace(name) == "" || (strings.TrimSpace(ip) == "" && strings.TrimSpace(mac) == "") {
		return false, "设备名称为空或者设备ip和mac地址均为空"
	}
	if ip != "" && !reIp.MatchString(ip) {
		return false, "ip地址格式不正确"
	}
	if mac != "" && !reMac.MatchString(mac) {
		return false, "mac地址格式不正确"
	}
	if wol && mac == "" {
		return false, "WOL需要mac地址"
	}
	if udp && ip == "" {
		return false, "UDP需要ip地址"
	}
	if submask < 1 || submask > 32 {
		return false, "子网掩码位数位于1~32之间"
	}
	return true, ""
}

func SetDevice(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		devices := getDevices()
		json.NewEncoder(w).Encode(&ApiReturn{
			Retcode: 0,
			Message: "OK",
			Data: &AllDevices{
				Count:   len(devices),
				Devices: devices,
			},
		})
	case "PUT":
		var body AllDevices
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			ApiErr(w)
			return
		}
		var devices []Device
		var msg string

		for _, dev := range body.Devices {
			valid, m := checkDeviceValid(dev.DeviceName, dev.DeviceIp, dev.DeviceMac, dev.DeviceUdp, dev.DeviceWol, dev.DeviceSubmask)
			if valid {
				devices = append(devices, dev)
			} else {
				msg = msg + m + " "
			}
		}
		if len(devices) == 0 {
			ApiErrMsg(w, msg+"No item to add")
			return
		}
		ok := addDevice(devices)
		if ok {
			ApiOkMsg(w, msg+"OK")
		} else {
			ApiErrMsg(w, msg+"请求错误")
		}
	case "POST":
		var body Device
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			ApiErr(w)
			return
		}
		if getDeviceById(body.DeviceId) == nil {
			ApiErrMsg(w, "设备不存在")
			return
		}
		valid, m := checkDeviceValid(body.DeviceName, body.DeviceIp, body.DeviceMac, body.DeviceUdp, body.DeviceWol, body.DeviceSubmask)
		if !valid {
			ApiErrMsg(w, m)
			return
		}
		ok := setDevice(body.DeviceId, body.DeviceName, body.DeviceIp, body.DeviceMac, body.DeviceUdp, body.DeviceWol, body.DeviceSubmask)
		if ok {
			ApiOk(w)
		} else {
			ApiErr(w)
		}
	case "DELETE":
		var body Device
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			ApiErr(w)
			return
		}
		if getDeviceById(body.DeviceId) == nil {
			ApiErrMsg(w, "设备不存在")
			return
		}
		ok := deleteDevice(body.DeviceId)
		if ok {
			ApiOk(w)
		} else {
			ApiErr(w)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
