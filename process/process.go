package process

import (
	"encoding/json"
	"os"
)

//var processListPath string = "process/processList.json"
//var schedulePath string = "process/schedule.json"

type Process struct {
	Name			string
	URIWhitelist    []string
	StartHourUTC    int
	StartMinuteUTC  int
	DurationHours   int
	DurationMinutes int
}

type Processes []Process

func InsertProcess(fileName string, process Process) error {

	var processes Processes
	readErr := ReadProcesses(fileName, &processes)

	if readErr != nil {
		return readErr
	}

	processes = append(processes, process)

	processByteData, marshalErr := json.Marshal(processes)

	if marshalErr != nil {
		return marshalErr
	}

	os.WriteFile(fileName, processByteData, os.ModePerm)

	return nil
}

func ReadProcesses(fileName string, processes *Processes) error {
	fileData, readErr := os.ReadFile(fileName)
	
	if readErr != nil {
		return readErr
	}

	_ = json.Unmarshal(fileData, processes)

	return nil
}

func RemoveProcess(process Process) {}

func UpdateProcess(process Process) {}

