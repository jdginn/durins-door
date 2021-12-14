package parser

import (
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
		t.Log("Error calling GetReader: %v", err)
		t.Fail()
	}
}

func TestHasAttr(t *testing.T) {
	reader, _ := GetReader(testcaseFilename)
}

func TestParseLocation(t *testing.T) {
}

func TestGetTypeDie(t *testing.T) {
}

func TestListAllAttributes(t *testing.T) {
}
