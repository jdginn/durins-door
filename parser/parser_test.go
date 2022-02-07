package parser

import (
	"debug/dwarf"
  "debug/macho"
	"github.com/stretchr/testify/assert"
	"testing"
)

// For now, this functionality all relies upon having a reader object. The only good
// idea I have for creating that reader is to read a DWARF file and the best way to do that
// is to simply compile a testcase. I am using this: https://github.com/jdginn/testcase-compiler

// The downside here is that these tests are hostage to changes in that testcase.
var testcaseFilename = "testcase-compiler/testcase.out.dSYM/Contents/Resources/DWARF/testcase.out"

func getReaderFromFile(fileName string) (*dwarf.Reader, error) {
  fh, err := macho.Open(fileName) 
  if err != nil {
    return nil, err
  }
  return GetReader(fh)
}

func TestgetReaderFromFile(t *testing.T) {
	// For now, just assume testcase is always located in the right place
	_, err := getReaderFromFile(testcaseFilename)
	assert.Nil(t, err)
}

func testGetEntry(t *testing.T, requestedName string) *dwarf.Entry {
	reader, _ := getReaderFromFile(testcaseFilename)
	entry, err := GetEntry(reader, requestedName)
	assert.Nil(t, err)
	assert.Equal(t, entry.Val(dwarf.AttrName), requestedName)
	return entry
}

func shouldFailGetEntry(t *testing.T, requestedName string, errorString string) {
	reader, _ := getReaderFromFile(testcaseFilename)
	_, err := GetEntry(reader, requestedName)
	// Test that we can read twice in a row without building a new reader
	_, err = GetEntry(reader, requestedName)
	assert.NotNil(t, err)
}

func TestGetEntry(t *testing.T) {
	testGetEntry(t, "formula_1_teams")
	testGetEntry(t, "testcase.cpp")
	testGetEntry(t, "Driver")
	testGetEntry(t, "Driver")
	shouldFailGetEntry(t, "badname", "entry could not be found")
}

func testGetTypeEntry(t *testing.T, reader *dwarf.Reader, entryName string) *dwarf.Entry {
	entry, err := GetEntry(reader, entryName)
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
