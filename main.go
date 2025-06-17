package main

import (
	"fmt"
	process "github.com/AndochBonin/myDaemon/process"
)

func main() {
	fmt.Println("hello, world!")
	insertErr := process.InsertProcess("process/processList.json", process.Process{Name: ""})
	fmt.Println(insertErr)
}
