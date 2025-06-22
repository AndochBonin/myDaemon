package main

import (
	"github.com/AndochBonin/myDaemon/process"
	"github.com/AndochBonin/myDaemon/program"
)

var programListFile string = "./program/programList.json"

func main() {
	testProgram := program.Program{Name: "test"}
	program.CreateProgram(programListFile, testProgram)

	schedule := process.GetScheduler()

	schedule.AddProcess(process.Process{Program: testProgram})
}
