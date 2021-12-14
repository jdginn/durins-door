package parser

import (
  "debug/dwarf"
	"testing"
)

// For now, this functionality all relies upon having a reader object. The only good
// idea I have for creating that reader is to read a DWARF file and the best way to do that
// is to simply compile a testcase. I am using this: https://github.com/jdginn/testcase-compiler

// The downside here is that these tests are hostage to changes in that testcase.
var testcaseFilename = "../../testcase-compiler/main.out.dSYM/Contents/Resources/DWARF/main.out"

func TestGetReader(t *testing.T) {
	// For now, just assume testcase is always located in the right place
	_, err := GetReader(testcaseFilename)
	if err != nil {
		t.Fatal("Error calling GetReader", err)
	}
}

func testGetEntry(t *testing.T, requestedName string) {
	reader, _ := GetReader(testcaseFilename)
  entry, err := GetEntry(reader, requestedName)
  if err != nil {
    t.Fatal(err)
  }
  foundName := entry.AttrField(dwarf.AttrName).Val
  if foundName != requestedName {
    t.Log("Found the wrong entry.")
    t.Log("  Requested entry: ", requestedName)
    t.Log("  Found entry: ", foundName)
    t.Fail()
  }
}

func shouldFailGetEntry(t *testing.T, requestedName string, errorString string) {
	reader, _ := GetReader(testcaseFilename)
  _, err := GetEntry(reader, requestedName)
  if err == nil {
    t.Fatal("Should have failed locating entry", requestedName)
  }
}

func TestGetEntry(t *testing.T) {
  testGetEntry(t, "formula_1_teams")
  testGetEntry(t, "main.cpp")
  testGetEntry(t, "drivers")
  testGetEntry(t, "Driver")
  shouldFailGetEntry(t, "badname", "entry could not be found")
}

func testGetTypeEntry(t *testing.T, entryName string, expectedTag dwarf.Tag) {
  reader, _ := GetReader(testcaseFilename)
  entry, err := GetEntry(reader, entryName)
  if err != nil {
    t.Fatal(err)
  }
  typeEntry, err := GetTypeEntry(reader, entry)
  if err != nil {
    t.Fatal("Could not get typedef", err)
  }
  if typeEntry.Tag != expectedTag {
    t.Fatal("Did not find the expected tag")
  }
}

func TestGetTypeEntry(t *testing.T) {
  testGetTypeEntry(t, "formula_1_teams", dwarf.TagArrayType)
  testGetTypeEntry(t, "drivers", dwarf.TagArrayType)
  testGetTypeEntry(t, "Driver", dwarf.TagStructType)
  testGetTypeEntry(t, "char", dwarf.TagBaseType)
}

func TestParseLocation(t *testing.T) {
}

func TestGetTypeDie(t *testing.T) {
}

func TestListAllAttributes(t *testing.T) {
}
