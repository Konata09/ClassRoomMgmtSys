package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"net"
	"time"
)

type DetectRes struct {
	id  int
	res int
}

var controllerQueryPayload = []byte("\x4c\x69\x67\x68\x74\x6f\x6e\xfe\x08\x14\x0a\x00\x26\xff")

func pingSingle(ip string, id int, c chan DetectRes) {
	var pingres DetectRes
	pingres.id = id
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		fmt.Printf("%s when ping to %s\n", err, ip)
		pingres.res = -1
		c <- pingres
		return
	}
	pinger.Timeout = 1 * time.Second
	pinger.Count = 1
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		fmt.Printf("%s when ping to %s\n", err, ip)
		pingres.res = -1
		c <- pingres
		return
	}
	stat := pinger.Statistics()
	var res int
	if stat.PacketLoss == 100 {
		res = -1
	} else {
		res = int(stat.AvgRtt.Microseconds())
	}
	pingres.res = res
	c <- pingres
}

func pingDevices(devices []Device, done chan int) {
	size := len(devices)
	c := make(chan DetectRes, size)
	for i, device := range devices {
		go pingSingle(device.DeviceIp, i, c)
	}
	for i := 0; i < size; i++ {
		pingres := <-c
		devices[pingres.id].pingRes = pingres.res
	}
	done <- 1
}

func getControllerStatusSingle(ip string, id int, c chan DetectRes) {
	var pingres DetectRes
	pingres.id = id

	pc, err := net.ListenUDP("udp4", nil)
	if err != nil {
		fmt.Printf("%s when getControollerStatus from %s\n", err, ip)
		pingres.res = -1
		c <- pingres
		return
	}
	defer pc.Close()
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", ip, 4001))
	if err != nil {
		fmt.Printf("%s when getControollerStatus from %s\n", err, ip)
		pingres.res = -1
		c <- pingres
		return
	}
	_, err = pc.WriteTo(controllerQueryPayload, addr)
	if err != nil {
		fmt.Printf("%s when getControollerStatus from %s\n", err, ip)
		pingres.res = -1
		c <- pingres
		return
	}
	buf := make([]byte, 8)
	pc.SetReadDeadline(time.Now().Add(time.Second * 1))
	_, _, err = pc.ReadFrom(buf)
	if err != nil {
		fmt.Printf("%s when getControollerStatus from %s\n", err, ip)
		pingres.res = -1
		c <- pingres
		return
	}
	/**
	res:
		-1 Fail
		1 ON
		2 OFF
		3 WAIT
		4 Unkonw
	*/
	if buf[4] == '\x12' && buf[5] == '\x70' {
		pingres.res = 1
	} else if buf[4] == '\x02' && buf[5] == '\x60' {
		pingres.res = 1
	} else if buf[4] == '\x02' && buf[5] == '\x70' {
		pingres.res = 1
	} else if buf[4] == '\x00' && buf[5] == '\x00' {
		pingres.res = 2
	} else if buf[4] == '\x00' && buf[5] == '\x05' {
		pingres.res = 2
	} else if buf[4] == '\x00' && buf[5] == '\x10' {
		pingres.res = 2
	} else if buf[4] == '\x00' && buf[5] == '\x60' {
		pingres.res = 2
	} else if buf[4] == '\x10' && buf[5] == '\x00' {
		pingres.res = 2
	} else if buf[4] == '\x03' && buf[5] == '\x20' {
		pingres.res = 3
	} else if buf[4] == '\x03' && buf[5] == '\x30' {
		pingres.res = 3
	} else if buf[4] == '\x03' && buf[5] == '\x60' {
		pingres.res = 3
	} else if buf[4] == '\x03' && buf[5] == '\x70' {
		pingres.res = 3
	} else {
		pingres.res = 4
	}
	c <- pingres
	return
}

func getControllersStatus(devices []Device, done chan int) {
	c := make(chan DetectRes, 50)
	var size int
	for i, device := range devices {
		if device.DeviceTypeId == 1 {
			go getControllerStatusSingle(device.DeviceIp, i, c)
			size++
		}
	}
	for i := 0; i < size; i++ {
		detectRes := <-c
		devices[detectRes.id].status = detectRes.res
	}
	done <- 2
}
