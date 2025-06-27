package program

import (
	"encoding/json"
	"os"
	"slices"
)

type Program struct {
	Name         string
	URIWhitelist []string
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

	return WriteProgramListToJSONFile(fileName, programList)
}

func ReadPrograms(fileName string, programList *ProgramList) error {
	fileData, readErr := os.ReadFile(fileName)

	if readErr != nil {
		if os.IsNotExist(readErr) {
			*programList = ProgramList{}
			return nil
		}
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

	if programID < 0 || programID >= len(programList) {
		return nil
	}
	programList = slices.Delete(programList, programID, programID+1)

	return WriteProgramListToJSONFile(fileName, programList)
}

func UpdateProgram(fileName string, programID int, program Program) error {
	var programList ProgramList
	readErr := ReadPrograms(fileName, &programList)

	if readErr != nil {
		return readErr
	}
	programList[programID] = program

	return WriteProgramListToJSONFile(fileName, programList)
}

func WriteProgramListToJSONFile(fileName string, programList ProgramList) error {
	programListByteData, marshalErr := json.MarshalIndent(programList, jsonPrefix, jsonIndent)

	if marshalErr != nil {
		return marshalErr
	}
	return os.WriteFile(fileName, programListByteData, 0644)
}
