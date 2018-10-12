package patch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Struct_Ignores_EmptyPatch(t *testing.T) {

	type Target struct {
		FirstName string
		LastName  string
		Salary    int
		Extra     string
	}
	var a = Target{
		FirstName: "Anakin",
		LastName:  "Skywalker",
		Salary:    123,
		Extra:     "unchanged",
	}
	var chg, err = Struct(&a, struct{}{})

	assert.NoError(t, err)
	assert.False(t, chg)
	assert.Equal(t, a.FirstName, "Anakin")
	assert.Equal(t, a.LastName, "Skywalker")
	assert.Equal(t, a.Salary, 123)
	assert.Equal(t, a.Extra, "unchanged")
}

func Test_Struct_Ignores_UnknownFields(t *testing.T) {

	type Target struct {
		FirstName string
		LastName  string
		Salary    int
	}
	type Patch struct {
		MiddleName string
		Perk       string
	}

	var a = Target{
		FirstName: "Anakin",
		LastName:  "Skywalker",
		Salary:    123,
	}
	var p = Patch{
		MiddleName: "Sheev",
		Perk:       "pilot",
	}

	var chg, err = Struct(&a, p)

	assert.NoError(t, err)
	assert.False(t, chg)
	assert.Equal(t, a.FirstName, "Anakin")
	assert.Equal(t, a.LastName, "Skywalker")
	assert.Equal(t, a.Salary, 123)
}

func Test_Struct_Ignores_UnexportedFields(t *testing.T) {

	type Target struct {
		Exported   string
		unexported string
	}
	type Patch struct {
		Exported   *string
		unexported *string
	}

	var a = Target{
		Exported:   "stormtrooper",
		unexported: "private",
	}
	var p = Patch{
		Exported:   strPtr("wookie"),
		unexported: strPtr("leutenant"),
	}

	var chg, err = Struct(&a, p)

	assert.NoError(t, err)
	assert.True(t, chg)
	assert.Equal(t, a.Exported, "wookie")
	assert.Equal(t, a.unexported, "private")
}

func Test_Struct_Ignores_ZeroValueFields(t *testing.T) {

	type Target struct {
		Name   string
		Salary int
		OnDuty bool
	}
	type Patch struct {
		Name   string
		Salary int
		OnDuty bool
	}

	var a = Target{
		Name:   "Han Solo",
		Salary: 15,
		OnDuty: true,
	}
	var p = Patch{
		Name:   "Chubacca",
		Salary: 0,
		OnDuty: false,
	}

	var chg, err = Struct(&a, p)

	assert.NoError(t, err)
	assert.True(t, chg)
	assert.Equal(t, a.Name, "Chubacca")
	assert.Equal(t, a.Salary, 15)
	assert.Equal(t, a.OnDuty, true)
}

func Test_Struct_Handles_Pointers(t *testing.T) {

	type Target struct {
		FirstName string
		LastName  string
		Salary    int
		OnDuty    bool
	}
	type Patch struct {
		FirstName *string
		LastName  *string
		Salary    *int
		OnDuty    *bool
	}

	var a = Target{
		FirstName: "Anakin",
		LastName:  "Skywalker",
		Salary:    123,
		OnDuty:    true,
	}
	var p = Patch{
		FirstName: strPtr("Darth"),
		LastName:  strPtr("Vader"),
		Salary:    intPtr(0),
		OnDuty:    boolPtr(false),
	}

	var chg, err = Struct(&a, p)

	assert.NoError(t, err)
	assert.True(t, chg)
	assert.Equal(t, a.FirstName, "Darth")
	assert.Equal(t, a.LastName, "Vader")
	assert.Equal(t, a.Salary, 0)
	assert.Equal(t, a.OnDuty, false)
}

func Test_Struct_Detects_WrongType(t *testing.T) {

	type Target struct {
		Name   string
		Salary int
	}
	type Patch struct {
		Name   string
		Salary float64 // wrong type
	}

	var a = Target{
		Name:   "Anakin Skywalker",
		Salary: 123,
	}
	var p = Patch{
		Name:   "Darth Vader",
		Salary: 500.0,
	}

	var _, err = Struct(&a, p)

	assert.Error(t, err)
	assert.Equal(t, a.Salary, 123)
}


// -- test helpers --

func strPtr(s string) *string {
	return &s
}
func intPtr(i int) *int {
	return &i
}
func boolPtr(b bool) *bool {
	return &b
}
