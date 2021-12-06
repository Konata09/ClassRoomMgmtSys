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
	pinger.Timeout = time.Second * 2
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
		logBoth("[ERR] %s when getControollerStatus ListenUDP from %s", err, ip)
		pingres.res = -1
		c <- pingres
		return
	}
	defer pc.Close()
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", ip, 4001))
	if err != nil {
		logBoth("[ERR] %s when getControollerStatus ResolveUDPAddr from %s", err, ip)
		pingres.res = -1
		c <- pingres
		return
	}
	_, err = pc.WriteTo(controllerQueryPayload, addr)
	if err != nil {
		logBoth("[ERR] %s when getControollerStatus WriteTo from %s", err, ip)
		pingres.res = -1
		c <- pingres
		return
	}
	buf := make([]byte, 8)
	err = pc.SetReadDeadline(time.Now().Add(time.Second * 1))
	if err != nil {
		fmt.Printf("SetReadDeadline Fail %s", err)
		return
	}
	_, _, err = pc.ReadFrom(buf)
	if err != nil {
		logBoth("[ERR] %s when getControollerStatus ReadFrom from %s", err, ip)
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

/**
  查询所有教室基本状态(中控 云盒 录直播)
*/
func getAllClassroomStatus() {
	classrooms := getClassrooms()
	done := make(chan int)
	lindges := getClassroomLindges()
	controllers := getClassroomControllers()
	go pingDevices(lindges, done)              // ping所有云盒 doneId:1
	go getControllersStatus(controllers, done) // 查询所有中控状态 doneId:2
	for i := 0; i < 2; i++ {
		doneId := <-done
		if doneId == 1 {
			for _, lindge := range lindges {
				for i, classroom := range classrooms {
					if classroom.Id == lindge.DeviceClassId {
						classrooms[i].Lindge = lindge.pingRes
						break
					}
				}
			}

		} else if doneId == 2 {
			for _, controller := range controllers {
				for i, classroom := range classrooms {
					if classroom.Id == controller.DeviceClassId {
						classrooms[i].Controller = controller.status
						break
					}
				}
			}
		}
	}
	//liveStatus := getLiveStatusFromDy() // 获取录直播状态
}

func handleStatusResults(doneId int, devices []Device) {

}
