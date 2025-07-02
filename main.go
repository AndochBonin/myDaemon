package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/AndochBonin/myDaemon/process"
	"github.com/AndochBonin/myDaemon/tui"
)

var exceptions = []string{"explorer", "Taskmgr", "WindowsTerminal", "TextInputHost"}

func killProcesses(whitelistMap map[string]bool) error {
	out, err := exec.Command("powershell", "-Command",
		`Get-Process | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object Name,Id | ConvertTo-Csv -NoTypeInformation`).Output()

	if err != nil {
		return err
	}

	processes, err := csv.NewReader(bytes.NewReader(out)).ReadAll()
	if err != nil {
		return err
	}

	for _, process := range processes[1:] {
		name := process[0]
		pid := process[1]
		if _, found := whitelistMap[strings.ToLower(name)]; !found {
			err := exec.Command("taskkill", "/PID", pid, "/F").Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func RunSchedule(scheduler *process.Scheduler) {
	var process *process.Process
	for {
		scheduler.UpdateSchedule()
		process = scheduler.GetCurrentProcess()
		if process != nil {
			whitelistMap := make(map[string]bool)
			for _, name := range append(process.Program.URIWhitelist, exceptions...) {
				whitelistMap[strings.ToLower(name)] = true
			}
			killProcesses(whitelistMap)			
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	scheduler := process.GetScheduler()
	go RunSchedule(scheduler)
	err := tui.Run()
	if err != nil {
		fmt.Println("welp oops")
	}
}
