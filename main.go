package main

import (
	"fmt"
	//"github.com/AndochBonin/myDaemon/process"
	//"github.com/AndochBonin/myDaemon/program"
	"github.com/AndochBonin/myDaemon/tui"
)

//var programListFile string = "./program/programList.json"

func main() {

	err := tui.Run()

	if err != nil {
		fmt.Println("welp oops")
	}
}
