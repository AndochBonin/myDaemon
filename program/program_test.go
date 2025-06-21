package program

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

const testFileName = "test_programs.json"

func cleanupTestFile(t *testing.T) {
	t.Helper()
	if err := os.Remove(testFileName); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to clean up test file: %v", err)
	}
}

func TestReadPrograms(t *testing.T) {
	cleanupTestFile(t)
	defer cleanupTestFile(t)

	expected := ProgramList{
		{
			Name:         "ReadTestApp",
			URIWhitelist: []string{"https://read.example.com"},
		},
	}

	data, _ := json.MarshalIndent(expected, "", "    ")
	err := os.WriteFile(testFileName, data, 0644)
	if err != nil {
		t.Fatalf("Failed to set up test file: %v", err)
	}

	var actual ProgramList
	err = ReadPrograms(testFileName, &actual)
	if err != nil {
		t.Fatalf("ReadPrograms failed: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestCreateProgram(t *testing.T) {
	cleanupTestFile(t)
	defer cleanupTestFile(t)

	prog := Program{
		Name:         "CreateTestApp",
		URIWhitelist: []string{"https://create.example.com"},
	}

	err := CreateProgram(testFileName, prog)
	if err != nil {
		t.Fatalf("CreateProgram failed: %v", err)
	}

	var progList ProgramList
	err = ReadPrograms(testFileName, &progList)
	if err != nil {
		t.Fatalf("ReadPrograms after create failed: %v", err)
	}

	if len(progList) != 1 || !reflect.DeepEqual(progList[0], prog) {
		t.Errorf("CreateProgram did not persist data correctly. Got: %v", progList)
	}
}

func TestUpdateProgram(t *testing.T) {
	cleanupTestFile(t)
	defer cleanupTestFile(t)

	original := Program{Name: "OriginalApp", URIWhitelist: []string{"https://orig.com"}}
	_ = CreateProgram(testFileName, original)

	updated := Program{Name: "UpdatedApp", URIWhitelist: []string{"https://updated.com"}}
	err := UpdateProgram(testFileName, 0, updated)
	if err != nil {
		t.Fatalf("UpdateProgram failed: %v", err)
	}

	var list ProgramList
	err = ReadPrograms(testFileName, &list)
	if err != nil {
		t.Fatalf("ReadPrograms after update failed: %v", err)
	}

	if len(list) != 1 || !reflect.DeepEqual(list[0], updated) {
		t.Errorf("Expected updated program %v, got %v", updated, list[0])
	}
}

func TestDeleteProgram(t *testing.T) {
	cleanupTestFile(t)
	defer cleanupTestFile(t)

	toDelete := Program{Name: "ToDelete", URIWhitelist: []string{}}
	toKeep := Program{Name: "ToKeep", URIWhitelist: []string{}}

	_ = CreateProgram(testFileName, toDelete)
	_ = CreateProgram(testFileName, toKeep)

	err := DeleteProgram(testFileName, 0)
	if err != nil {
		t.Fatalf("DeleteProgram failed: %v", err)
	}

	var list ProgramList
	err = ReadPrograms(testFileName, &list)
	if err != nil {
		t.Fatalf("ReadPrograms after delete failed: %v", err)
	}

	if len(list) != 1 || !reflect.DeepEqual(list[0], toKeep) {
		t.Errorf("Expected one program %v, got: %v", toKeep, list)
	}
}
