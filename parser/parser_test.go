package parser

import (
	"debug/dwarf"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jdginn/durins-door/explorer/plat"
)

// For now, this functionality all relies upon having a reader object. The only good
// idea I have for creating that reader is to read a DWARF file and the best way to do that
// is to simply compile a testcase. I am using this: https://github.com/jdginn/testcase-compiler
//
// The downside here is that these tests are hostage to changes in that testcase.
var testcaseFilename = "../testcase-compiler/testcase.dwarf"

func getReaderFromFile(fileName string) (*dwarf.Reader, error) {
	fh, err := plat.GetReaderFromFile(fileName)
	if err != nil {
		wd, _ := os.Getwd()
		panic(fmt.Errorf("Could not open file %s:\n\ncwd:  %s\n\nerror message:\n\t%s", fileName, wd, err))
	}
	return GetReader(fh)
}

func TestGetReaderFromFile(t *testing.T) {
	// For now, just assume testcase is always located in the right place
	_, err := getReaderFromFile(testcaseFilename)
	assert.NoError(t, err)
}

func testGetEntry(t *testing.T, requestedName string) *dwarf.Entry {
	reader, _ := getReaderFromFile(testcaseFilename)
	entry, _, err := GetEntry(reader, requestedName)
	assert.NoError(t, err)
	assert.Equal(t, entry.Val(dwarf.AttrName), requestedName)
	return entry
}

func shouldFailGetEntry(t *testing.T, requestedName string, errorString string) {
	reader, _ := getReaderFromFile(testcaseFilename)
	_, _, err := GetEntry(reader, requestedName)
	assert.Error(t, err)
}

func TestGetEntry(t *testing.T) {
	testGetEntry(t, "formula_1_teams")
	testGetEntry(t, "testcase.cpp")
	testGetEntry(t, "Driver")
	testGetEntry(t, "Driver")
	shouldFailGetEntry(t, "badname", "entry could not be found")
}

func testGetTypeEntry(t *testing.T, reader *dwarf.Reader, entryName string) *dwarf.Entry {
	entry, _, err := GetEntry(reader, entryName)
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatal("Failed to retrieve entry " + entryName)
	}
	typeEntry, err := GetTypeEntry(reader, entry)
	if err != nil {
		t.Fatal("Could not get typedef", err)
	}
	if typeEntry == nil {
		t.Fatal("Failed to retrieve type entry for " + entryName)
	}
	return typeEntry
}

func TestGetTypeEntry(t *testing.T) {
	reader, _ := getReaderFromFile(testcaseFilename)
	var e *dwarf.Entry
	e = testGetTypeEntry(t, reader, "formula_1_teams")
	if e == nil {
		t.Fatal("entry is nil")
	}
	assert.Equal(t, e.Tag, dwarf.TagArrayType)
	e = testGetTypeEntry(t, reader, "drivers")
	assert.Equal(t, e.Tag, dwarf.TagArrayType)
	e = testGetTypeEntry(t, reader, "Driver")
	assert.Equal(t, e.Tag, dwarf.TagStructType)
	assert.Equal(t, e.Val(dwarf.AttrByteSize), int64(12))
	e = testGetTypeEntry(t, reader, "char")
	assert.Equal(t, e.Tag, dwarf.TagBaseType)
	assert.Equal(t, e.Val(dwarf.AttrByteSize), int64(1))
	// Make sure we can get entries we've already gotten
	e = testGetTypeEntry(t, reader, "formula_1_teams")
	assert.Equal(t, e.Tag, dwarf.TagArrayType)
	e = testGetTypeEntry(t, reader, "drivers")
	assert.Equal(t, e.Tag, dwarf.TagArrayType)
}

func TestParseLocation(t *testing.T) {
}

func TestGetTypeDie(t *testing.T) {
}

func TestListAllAttributes(t *testing.T) {
}

func TestGetEntriesOnLevel(t *testing.T) {
	reader, _ := getReaderFromFile(testcaseFilename)
	entries, err := GetChildren(reader, func(*dwarf.Entry) bool {
    return true
  })
	assert.NoError(t, err)
	fmt.Println(entries)
	assert.Equal(t, 1, len(entries))

  _, _, err = GetEntry(reader, "formula_1_teams")
	entries, err = GetChildren(reader, func(e *dwarf.Entry) bool {
    return e.Tag == dwarf.TagVariable
  })
	assert.NoError(t, err)
  fmt.Printf("Entries:\n")
  for _, e := range entries {
    fmt.Printf("\t%s\n", e.Val(dwarf.AttrName))
  }

  _, _, err = GetEntry(reader, "mercedes")
	entries, err = GetChildren(reader, func(e *dwarf.Entry) bool {
    return e.Tag == dwarf.TagVariable
  })
	assert.NoError(t, err)
  fmt.Printf("Entries:\n")
  for _, e := range entries {
    fmt.Printf("\t%s\n", e.Val(dwarf.AttrName))
  }
}

func TestGetCUs(t *testing.T) {
	reader, _ := getReaderFromFile(testcaseFilename)
	entries, err := GetCUs(reader)
	assert.NoError(t, err)
	fmt.Println(entries)
	assert.Equal(t, 1, len(entries))
}
