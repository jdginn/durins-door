package file

import(
  "debug/macho"
  "os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testcaseDwarfFile = "../testcase-compiler/testcase.dwarf"
var testcaseBinFile = "../testcase-compiler/testcase.out"

func TestNewVariableWrapper(t *testing.T) {
  dbgReader, err := macho.Open(testcaseDwarfFile)
  assert.NoError(t, err)
  binReader, err := os.Open(testcaseBinFile)
  assert.NoError(t, err)
  proxy, err := NewVariableWrapper(dbgReader, binReader, "formula_1_teams")
  assert.NoError(t, err)
  assert.NotNil(t, proxy)
}

// func TestReadVariableWrapper(t *testing.T) {
//   dbgReader, err := macho.Open(testcaseDwarfFile)
//   binReader, err := os.Open(testcaseBinFile)
//   proxy, err := NewVariableWrapper(dbgReader, binReader, "formula_1_teams")
//   proxy.Read()
//   assert.Equal(t, , proxy.Get())
//   assert.Equal(t, , proxy.GetField())
// }

// func TestWriteVariableWrapper(t *testing.T) {
//   dbgReader, err := macho.Open(testcaseDwarfFile)
//   binReader, err := os.Open(testcaseBinFile)
//   proxy, err := NewVariableWrapper(dbgReader, binReader, "formula_1_teams")
//   proxy.SetField()
//   proxy.SetField()
//   proxy.SetField()
//   proxy.Write()
//   assert.Equal(t, proxy.Read())
// }
