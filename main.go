package main

import (
	"fmt"
	"time"

	"github.com/AndochBonin/myDaemon/process"
	"github.com/AndochBonin/myDaemon/program"
	"github.com/AndochBonin/myDaemon/tui"
)

var programListFile string = "./program/programList.json"

func main() {
	testProgram := program.Program{Name: "aha a new one"}
	program.CreateProgram(programListFile, testProgram)

	schedule := process.GetScheduler()

	schedule.AddProcess(process.Process{Program: testProgram, StartTime: time.Now(), EndTime: time.Now().Add(time.Minute * 2)})

	err := tui.Run()

	if err != nil {
		fmt.Println("welp oops")
	}
}
