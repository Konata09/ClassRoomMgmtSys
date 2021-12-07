package main

import "time"
import "github.com/go-co-op/gocron"

func updateCron() {
	s := gocron.NewScheduler(time.Local)
	s.Every(10).Seconds().Do(fetchAllClassroomStatus)
	s.Every(3).Minutes().Do(fetchAllClassroomDeviceStatus)
	s.SingletonMode()
	s.StartAsync()
}
