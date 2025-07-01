package main

import (
	"fmt"
	//"time"

	"github.com/AndochBonin/myDaemon/process"

	//"github.com/AndochBonin/myDaemon/program"
	"github.com/AndochBonin/myDaemon/tui"
)

//var programListFile string = "./program/programList.json"
func UpdateSchedule(scheduler *process.Scheduler) {
	for {
		scheduler.UpdateSchedule()
		//time.Sleep(time.Second)
	}
}
 
func main() {
	scheduler := process.GetScheduler()
	go UpdateSchedule(scheduler)
	err := tui.Run()
	if err != nil {
		fmt.Println("welp oops")
	}
}
