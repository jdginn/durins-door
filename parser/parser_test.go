package parser

import "testing"

// For now, just assume testcase is always located in the right place

func TestGetReader(t *testing.T) {
  testcaseFilename := "../../testcase-compiler/main.out.dSYM/Contents/Resources/DWARF/main.out"
  _, err := GetReader(testcaseFilename)
  if err != nil {
    t.Fail()
  }
}

func TestHasAttr(t *testing.T) {
}

func TestParseLocation(t *testing.T) {
}

func TestGetTypeDie(t *testing.T) {
}

func TestListAllAttributes(t *testing.T) {
}
