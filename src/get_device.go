package main

import (
	"encoding/json"
	"net/http"
)

type UserDevice struct {
	DeviceId   int    `json:"id"`
	DeviceName string `json:"name"`
	DeviceUdp  bool   `json:"udp"`
	DeviceWol  bool   `json:"wol"`
}

type DeviceList struct {
	Count   int          `json:"count"`
	Devices []UserDevice `json:"devices"`
}

func GetDevice(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		devices := getUserDevices()
		json.NewEncoder(w).Encode(&ApiReturn{
			Retcode: 0,
			Message: "OK",
			Data: &DeviceList{
				Count:   len(devices),
				Devices: devices,
			},
		})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
