package patch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Struct_Works(t *testing.T) {

	type Target struct {
		FirstName string
		LastName  string
		Salary    int
		Extra     string
	}
	type Patch struct {
		FirstName  *string
		LastName   *string
		Salary     int
		unexported bool
		NonExisting []byte
	}

	var a = Target{
		FirstName: "Anakin",
		LastName:  "Skywalker",
		Salary:    123,
		Extra:     "unchanged",
	}
	var chg, err = Struct(&a, Patch{})

	assert.NoError(t, err)
	assert.False(t, chg)
	assert.Equal(t, a.FirstName, "Anakin")
	assert.Equal(t, a.LastName, "Skywalker")
	assert.Equal(t, a.Salary, 123)
	assert.Equal(t, a.Extra, "unchanged")

	var p = Patch{
		FirstName:  strPtr("Darth"),
		LastName:   strPtr("Vader"),
		unexported: true,
        NonExisting: []byte{10,20},
	}
	chg, err = Struct(&a, p)

	assert.NoError(t, err)
	assert.True(t, chg)
	assert.Equal(t, a.FirstName, "Darth")
	assert.Equal(t, a.LastName, "Vader")
	assert.Equal(t, a.Salary, 123)
	assert.Equal(t, a.Extra, "unchanged")
}

func Test_Struct_ThowsError(t *testing.T) {

	type Target struct {
		Name   string
		Salary int
	}
	type Patch struct {
		Name   *string
		Salary float64  // wrong type
	}

	var a = Target{
		Name:   "Anakin Skywalker",
		Salary: 123,
	}
	var chg, err = Struct(&a, Patch{})

	assert.NoError(t, err)
	assert.False(t, chg)
	assert.Equal(t, a.Name, "Anakin Skywalker")
	assert.Equal(t, a.Salary, 123)

	var p = Patch{Name: strPtr("Darth Vader"), Salary: 500.0}
	chg, err = Struct(&a, p)

	assert.Error(t, err)
	assert.Equal(t, a.Salary, 123)
}

func strPtr(s string) *string {
	return &s
}
