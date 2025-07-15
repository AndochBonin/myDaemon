package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"github.com/AndochBonin/myDaemon/process"
	"github.com/AndochBonin/myDaemon/tui"
)

var exePath, _ = os.Executable()
var scheduleFile string = filepath.Join(filepath.Dir(exePath), "storage", "schedule.json")

var currentProcess *process.Process
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
	for {
		scheduler.UpdateSchedule(scheduleFile)
		currentProcess = scheduler.GetCurrentProcess()
		if currentProcess != nil {
			whitelistMap := make(map[string]bool)
			for _, name := range append(currentProcess.Program.AppWhitelist, exceptions...) {
				whitelistMap[strings.ToLower(name)] = true
			}
			killProcesses(whitelistMap)
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	scheduler := process.GetScheduler()
	process.ReadScheduleFromFile(scheduleFile, &scheduler.Schedule)
	go RunSchedule(scheduler)
	certFile := "C:\\Certs\\mydaemon.pem"
	keyFile := "C:\\Certs\\mydaemon.key"
	if !isKeyCertPairExist(certFile, keyFile) {
		log.Fatal("Key Cert pair does not exist")
	}
	go RunTLSProxy(certFile, keyFile)
	err := tui.Run()
	if err != nil {
		fmt.Println("welp oops")
	}
}

func isKeyCertPairExist(certFile string, keyFile string) bool {
	_, certErr := os.Stat(certFile)
	_, keyErr := os.Stat(keyFile)
	if os.IsNotExist(certErr) || os.IsNotExist(keyErr) {
		return false
	}
	return true
}
