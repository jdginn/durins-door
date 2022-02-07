package file

import(
  "debug/macho"
  "os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testcaseDwarfFile = "../../testcase-compiler/main.out.dSYM/Contents/Resources/DWARF/main.out"
var testcaseBinFile = "../../testcase-compiler/main.out"

func TestNewVariableWrapper(t *testing.T) {
  dbgReader, err := macho.Open(testcaseDwarfFile)
  assert.Nil(t, err)
  binReader, err := os.Open(testcaseBinFile)
  assert.Nil(t, err)
  proxy, err := NewVariableWrapper(dbgReader, binReader, "formula_1_teams")
  assert.Nil(t, err)
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
