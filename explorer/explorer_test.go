package explorer_test

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jdginn/durins-door/explorer"
)

var testcaseFilename = "../testcase-compiler/testcase.dwarf"

func TestNewExplorer(t *testing.T) {
	ex := explorer.NewExplorer()
	assert.Equal(t, []string(nil), ex.ListChildren())
	assert.Equal(t, "modeCUs", ex.CurrMode())
	assert.Equal(t, "all CUs", ex.CurrName())

	assert.Error(t, ex.CreateReaderFromFile("invalid_file"))
	assert.NoError(t, ex.CreateReaderFromFile(testcaseFilename))

	assert.NotPanics(t, func() { explorer.NewExplorerFromFile(testcaseFilename) })
}

func TestExplore(t *testing.T) {
	ex := explorer.NewExplorerFromFile(testcaseFilename)
	names, err := ex.ListCUs()
	assert.NoError(t, err)
	assert.Equal(t, []string{"testcase.cpp"}, names)

	assert.Error(t, ex.StepIntoChild("bad name"))

	err = ex.StepIntoChild("testcase.cpp")
	assert.NoError(t, err)
	assert.Equal(t, "modeEntry", ex.CurrMode())
	assert.Equal(t, "testcase.cpp", ex.CurrName())
	// variables := []string{
	// 	"formula_1_teams",
	// 	"red_bull",
	// 	"verstappen",
	// 	"perez",
	// 	"mercedes",
	// 	"hamilton",
	// 	"bottas",
	// }
	// assert.Equal(t, variables, ex.ListChildren())
	assert.Contains(t, ex.ListChildren(), "formula_1_teams")

	err = ex.StepIntoChild("formula_1_teams")
	assert.NoError(t, err)
	assert.Equal(t, "modeProxy", ex.CurrMode())
	assert.Equal(t, "formula_1_teams", ex.CurrName())
	members := []string{
		"drivers",
		"sponsors",
		"has_won_wdc",
		"last_wdc",
		"has_won_wcc",
		"last_wcc",
	}
	assert.Equal(t, members, ex.ListChildren())

	// TODO: Back doesn't work properly yet
	// assert.NoError(t, ex.Back())
	// assert.Equal(t, "modeEntry", ex.CurrMode())
	// assert.Equal(t, "testcase.cpp", ex.CurrName())
	//  fmt.Printf("\n\n\n%v\n\n\n", ex.ListChildren())
	//  assert.Contains(t, ex.ListChildren(), "formula_1_teams")

	//  err = ex.StepIntoChild("fomula_1_teams")
	//  assert.NoError(t, err)
	err = ex.StepIntoChild("drivers")
	assert.NoError(t, err)
	assert.Equal(t, "modeProxy", ex.CurrMode())
	assert.Equal(t, "drivers", ex.CurrName())
	assert.Equal(t, []string{"initials", "car_number", "has_won_wdc"}, ex.ListChildren())

	err = ex.GetType()
	assert.NoError(t, err)
	assert.Equal(t, "modeProxy", ex.CurrMode())
	assert.Equal(t, "drivers", ex.CurrName())
	assert.Equal(t, []string{"initials", "car_number", "has_won_wdc"}, ex.ListChildren())
}

func TestReadCUs(t *testing.T) {
	ex := explorer.NewExplorer()
	ex.CreateReaderFromFile(testcaseFilename)
	cus, err := ex.ListCUs()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cus))
}
