package process

import (
	"encoding/json"
	"os"
	"slices"
)

type Process struct {
	Name			string
	URIWhitelist    []string
	StartHourUTC    int
	StartMinuteUTC  int
	DurationHours   int
	DurationMinutes int
}

type Processes []Process

var jsonPrefix = ""
var jsonIndent = "    "

func InsertProcess(fileName string, process Process) error {

	var processes Processes
	readErr := ReadProcesses(fileName, &processes)

	if readErr != nil {
		return readErr
	}

	processes = append(processes, process)

	processByteData, marshalErr := json.MarshalIndent(processes, jsonPrefix, jsonIndent)

	if marshalErr != nil {
		return marshalErr
	}
	return os.WriteFile(fileName, processByteData, 0644)
}

func ReadProcesses(fileName string, processes *Processes) error {
	fileData, readErr := os.ReadFile(fileName)
	
	if readErr != nil {
		return readErr
	}

	_ = json.Unmarshal(fileData, processes)

	return nil
}

func RemoveProcess(fileName string, processID int) error {
	var processes Processes
	readErr := ReadProcesses(fileName, &processes)
	
	if readErr != nil {
		return readErr
	}

	processes = slices.Delete(processes, processID, processID + 1)

	processesByteData, marshalErr := json.MarshalIndent(processes, jsonPrefix, jsonIndent)

	if marshalErr != nil {
		return marshalErr
	}

	return os.WriteFile(fileName, processesByteData, 0644)
}

func UpdateProcess(fileName string, processID int, process Process) error {
	var processes Processes
	readErr := ReadProcesses(fileName, &processes)
	
	if readErr != nil {
		return readErr
	}

	processes[processID] = process

	processesByteData, marshalErr := json.MarshalIndent(processes, jsonPrefix, jsonIndent)

	if marshalErr != nil {
		return marshalErr
	}

	return os.WriteFile(fileName, processesByteData, 0644)
}

