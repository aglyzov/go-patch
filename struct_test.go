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
	assert.Equal(t, "Anakin", a.FirstName)
	assert.Equal(t, "Skywalker", a.LastName)
	assert.Equal(t, 123, a.Salary)
	assert.Equal(t, "unchanged", a.Extra)
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
	assert.Equal(t, "Anakin", a.FirstName)   // unchanged
	assert.Equal(t, "Skywalker", a.LastName) // unchanged
	assert.Equal(t, 123, a.Salary)           // unchanged
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
	assert.Equal(t, "wookie", a.Exported)
	assert.Equal(t, "private", a.unexported)
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
	assert.Equal(t, "Chubacca", a.Name)
	assert.Equal(t, 15, a.Salary)   // unchanged
	assert.Equal(t, true, a.OnDuty) // unchanged
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
	assert.Equal(t, "Darth", a.FirstName)
	assert.Equal(t, "Vader", a.LastName)
	assert.Equal(t, 0, a.Salary)      // changed despite it's a zero
	assert.Equal(t, false, a.OnDuty)  // changed despite it's a zero
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
	assert.Equal(t, 123, a.Salary) // unchanged
}

func Test_Struct_Patches_EmbeddedWithFlat(t *testing.T) {

	type Name struct {
		FirstName string
		LastName  string
	}
	type Target struct {
		Name   // embedded
		Salary int
	}
	type Patch struct {
		FirstName string
		LastName  string
		Salary    int
	}

	var a = Target{
		Name: Name{
			FirstName: "Anakin",
			LastName:  "Skywalker",
		},
		Salary: 123,
	}
	var p = Patch{
		FirstName: "Darth",
		LastName:  "Vader",
		Salary:    100500,
	}

	var chg, err = Struct(&a, p)

	assert.NoError(t, err)
	assert.True(t, chg)
	assert.Equal(t, "Darth", a.FirstName)
	assert.Equal(t, "Vader", a.LastName)
	assert.Equal(t, 100500, a.Salary)
}

func Test_Struct_Patches_FlatWithEmbedded(t *testing.T) {

	type Target struct {
		FirstName string
		LastName  string
		Salary    int
	}
	type Name struct {
		FirstName string
		LastName  string
	}
	type Patch struct {
		Name   // embedded
		Salary int
	}

	var a = Target{
        FirstName: "Anakin",
        LastName:  "Skywalker",
		Salary: 123,
	}
	var p = Patch{
		Name: Name{
            FirstName: "Darth",
            LastName:  "Vader",
		},
		Salary:    100500,
	}

	var chg, err = Struct(&a, p)

	assert.NoError(t, err)
	assert.True(t, chg)
	assert.Equal(t, "Darth", a.FirstName)
	assert.Equal(t, "Vader", a.LastName)
	assert.Equal(t, 100500, a.Salary)
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
