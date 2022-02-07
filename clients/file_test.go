package file

import(
  "debug/macho"
  "os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testcaseDwarfFile = "../../testcase-compiler/main.out.dSYM/Contents/Resources/DWARF/main.out"
var testcaseBinFile = "../../testcase-compiler/main.out"

func TestNewVariableProxy(t *testing.T) {
  dbgReader, err := macho.Open(testcaseDwarfFile)
  assert.Nil(t, err)
  binReader, err := os.Open(testcaseBinFile)
  assert.Nil(t, err)
  proxy, err := NewVariableProxy(dbgReader, binReader, "formula_1_teams")
  assert.Nil(t, err)
  assert.NotNil(t, proxy)
}
