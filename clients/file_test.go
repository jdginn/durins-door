package file

import(
  "os"
	"testing"

	"github.com/stretchr/testify/assert"
  "github.com/jdginn/dwarf-explore"
)

var testcaseDwarfFile = "../testcase-compiler/main.out.dSYM/Contents/Resources/DWARF/main.out"
var testcaseBinFile = "../testcase-compiler/main.out"

func TestReader(t *testing.T) {
	dwarfReader, err := parser.GetReader(testcaseDwarfFile)
  assert.Nil(t, err)
  binReader, err := os.Open(testcaseBinFile)
  assert.Nil(t, err)


}
