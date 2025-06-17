package main

import (
	"fmt"
	process "github.com/AndochBonin/myDaemon/process"
)

func main() {
	err := process.RemoveProcess("process/processList.json", 0)

	if err != nil {
		fmt.Println(err.Error())
	}
}
