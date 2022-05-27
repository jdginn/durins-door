package explorer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jdginn/durins-door/explorer"
)

var testcaseFilename = "../testcase-compiler/testcase.dwarf"

func TestReadCUs(t *testing.T) {
	ex := explorer.NewExplorer()
	ex.CreateReaderFromFile(testcaseFilename)
	cus, err := ex.ListCUs()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cus))
}
