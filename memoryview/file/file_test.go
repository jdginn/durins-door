package file_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jdginn/durins-door/memoryview"
	"github.com/jdginn/durins-door/memoryview/file"
)

var testcaseDwarfFile = "../../testcase-compiler/testcase.dwarf"
var testcaseBinFile = "../../testcase-compiler/testcase.out"

func wantsClient(c memoryview.MemoryView) {}

func TestInterfaceMembership(t *testing.T) {
	dummy, err := file.NewFromPath(testcaseBinFile)
	assert.NoError(t, err)

	wantsClient(dummy)
}

func TestReadWrite(t *testing.T) {
	f, err := os.Create("test.bin")
	defer os.Remove("test.bin")
	assert.NoError(t, err)
	c, err := file.New(f)
	assert.NoError(t, err)
	err = c.Write(0, []byte("\xfe\xed\xbe\xef"))
	assert.NoError(t, err)
	rdata, err := c.Read(0, 2)
	assert.Equal(t, []byte("\xfe\xed"), rdata)
	assert.NoError(t, err)
	rdata, err = c.Read(2, 2)
	assert.Equal(t, []byte("\xbe\xef"), rdata)
	assert.NoError(t, err)

	// Not enough space for this read
	rdata, err = c.Read(2, 10)
	assert.Error(t, err)

	err = c.Write(4, []byte("\x00\x00\x00\x00\x00\x00\x00\x00"))
	assert.NoError(t, err)

	// This should work now that we have padded out the file sufficiently
	rdata, err = c.Read(0, 12)
	assert.NoError(t, err)
	assert.Equal(t, []byte("\xfe\xed\xbe\xef\x00\x00\x00\x00\x00\x00\x00\x00"), rdata)

	c.Write(1, []byte("\x00\x00"))
	rdata, err = c.Read(0, 2)
	assert.Equal(t, []byte("\xfe\x00"), rdata)
}
