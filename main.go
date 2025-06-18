package main

import (
"github.com/AndochBonin/myDaemon/program"
)

var programListFile string = "./program/programList.json"

func main() {
	program.CreateProgram(programListFile, program.Program{})
}
