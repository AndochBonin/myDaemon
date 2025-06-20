package main

import (
"github.com/AndochBonin/myDaemon/program"
"github.com/AndochBonin/myDaemon/process"
)

var programListFile string = "./program/programList.json"

func main() {
	testProgram := program.Program{Name: "test"}
	program.CreateProgram(programListFile, testProgram)

	scheduler := process.GetScheduler()

	scheduler.AddProcess(process.Process{Program: testProgram})
}
