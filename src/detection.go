package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"time"
)

type pingRes struct {
	id  int
	res int
}

func pingSingle(ip string, id int, c chan pingRes) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		panic(err)
	}
	pinger.Timeout = 2 * time.Second
	pinger.Count = 1
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		panic(err)
	}
	stat := pinger.Statistics()
	var res int
	if stat.PacketLoss == 100 {
		res = -1
	} else {
		res = int(stat.AvgRtt.Microseconds())
	}
	fmt.Printf("ip: %s ping: %d\n",ip,res)
	var pingres pingRes
	pingres.id = id
	pingres.res = res
	c <- pingres
}

func pingDevices(devices []Device) {
	size := len(devices)
	c := make(chan pingRes, size)
	for i, device := range devices {
		go pingSingle(device.DeviceIp, i, c)
	}
	for i := 0; i < size; i++ {
		pingres := <- c
		devices[pingres.id].pingRes = pingres.res
	}
	return
}
