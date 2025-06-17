package program

import (
	"encoding/json"
	"os"
	"slices"
)

type Program struct {
	Name			string
	URIWhitelist    []string
}

type ProgramList []Program

var jsonPrefix = ""
var jsonIndent = "    "

func CreateProgram(fileName string, program Program) error {

	var programList ProgramList
	readErr := ReadPrograms(fileName, &programList)

	if readErr != nil {
		return readErr
	}

	programList = append(programList, program)

	programListByteData, marshalErr := json.MarshalIndent(programList, jsonPrefix, jsonIndent)

	if marshalErr != nil {
		return marshalErr
	}
	return os.WriteFile(fileName, programListByteData, 0644)
}

func ReadPrograms(fileName string, programList *ProgramList) error {
	fileData, readErr := os.ReadFile(fileName)
	
	if readErr != nil {
		return readErr
	}

	_ = json.Unmarshal(fileData, programList)

	return nil
}

func DeleteProgram(fileName string, programID int) error {
	var programList ProgramList
	readErr := ReadPrograms(fileName, &programList)
	
	if readErr != nil {
		return readErr
	}

	programList = slices.Delete(programList, programID, programID + 1)

	programListByteData, marshalErr := json.MarshalIndent(programList, jsonPrefix, jsonIndent)

	if marshalErr != nil {
		return marshalErr
	}

	return os.WriteFile(fileName, programListByteData, 0644)
}

func UpdateProgram(fileName string, programID int, program Program) error {
	var programList ProgramList
	readErr := ReadPrograms(fileName, &programList)
	
	if readErr != nil {
		return readErr
	}

	programList[programID] = program

	programListByteData, marshalErr := json.MarshalIndent(programList, jsonPrefix, jsonIndent)

	if marshalErr != nil {
		return marshalErr
	}

	return os.WriteFile(fileName, programListByteData, 0644)
}
