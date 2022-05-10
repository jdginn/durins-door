package file

import(
	"testing"

	"github.com/stretchr/testify/assert"

	// "github.com/jdginn/dwarf-explore"
)

var testcaseDwarfFile = "../testcase-compiler/testcase.dwarf"
var testcaseBinFile = "../testcase-compiler/testcase.out"

// func wantsClient(c Client) {}

func TestInterfaceMembership(t *testing.T) {
  _, err := NewFileClient(testcaseBinFile)
  assert.NoError(t, err)

  // wantsClient(dummy)
}
