package patch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Apply_Ignores_EmptyPatch(t *testing.T) {

	type Target struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Salary    int    `json:"salary"`
		Extra     string `json:"extra"`
	}
	var a = Target{
		FirstName: "Anakin",
		LastName:  "Skywalker",
		Salary:    123,
		Extra:     "unchanged",
	}
	var chg, err = Apply(&a, map[string]interface{}{})

	assert.NoError(t, err)
	assert.False(t, chg)
	assert.Equal(t, "Anakin", a.FirstName)
	assert.Equal(t, "Skywalker", a.LastName)
	assert.Equal(t, 123, a.Salary)
	assert.Equal(t, "unchanged", a.Extra)
}

func Test_Apply_Ignores_UnknownFields(t *testing.T) {

	type Target struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Salary    int    `json:"salary"`
	}

	var a = Target{
		FirstName: "Anakin",
		LastName:  "Skywalker",
		Salary:    123,
	}

	data := `{"middle_name":"Sheev", "perk": "pilot"}`
	p := make(map[string]interface{})
	json.Unmarshal([]byte(data), &p)

	var chg, err = Apply(&a, p)

	assert.NoError(t, err)
	assert.False(t, chg)
	assert.Equal(t, "Anakin", a.FirstName)   // unchanged
	assert.Equal(t, "Skywalker", a.LastName) // unchanged
	assert.Equal(t, 123, a.Salary)           // unchanged
}

func Test_Apply_Ignores_UnexportedFields(t *testing.T) {

	type Target struct {
		Exported   string `json:"exported,omitempty"`
		unexported string
	}
	var a = Target{
		Exported:   "stormtrooper",
		unexported: "private",
	}

	data := `{"exported":"wookie", "unexported": "leutenant"}`
	p := make(map[string]interface{})
	json.Unmarshal([]byte(data), &p)

	var chg, err = Apply(&a, p)

	assert.NoError(t, err)
	assert.True(t, chg)
	assert.Equal(t, "wookie", a.Exported)
	assert.Equal(t, "private", a.unexported)
}

func Test_Apply_ZeroValueFields(t *testing.T) {

	type Target struct {
		Name   string `json:"name"`
		Salary int64  `json:"salary"`
		OnDuty bool   `json:"on_duty"`
	}

	var a = Target{
		Name:   "Han Solo",
		Salary: 15,
		OnDuty: true,
	}

	data := `{"name":"Chubacca", "salary": 0, "on_duty": false}`
	p := make(map[string]interface{})
	json.Unmarshal([]byte(data), &p)

	var chg, err = Apply(&a, p)

	assert.NoError(t, err)
	assert.True(t, chg)
	assert.Equal(t, "Chubacca", a.Name)
	assert.Equal(t, int64(0), a.Salary)
	assert.Equal(t, false, a.OnDuty)
}

func Test_Apply_Detects_WrongType(t *testing.T) {

	type Target struct {
		Name   string `json:"name"`
		Salary int    `json:"salary"`
	}
	var a = Target{
		Name:   "Anakin Skywalker",
		Salary: 123,
	}

	data := `{"name":"Darth Vader", "salary": "euros"}`
	p := make(map[string]interface{})
	json.Unmarshal([]byte(data), &p)

	var _, err = Apply(&a, p)

	assert.Error(t, err)
	assert.Equal(t, 123, a.Salary) // unchanged
}

func Test_Apply_Sub_Structs(t *testing.T) {
	type TargetPerson struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	type Target struct {
		Contact *TargetPerson `json:"contact"`
		Salary  int           `json:"salary"`
	}

	var a = Target{
		Contact: &TargetPerson{
			FirstName: "Anakin",
			LastName:  "Skywalker",
		},
		Salary: 123,
	}

	data := `{"contact": {"first_name":"Darth", "last_name": "Vader"}, "salary": 100500}`
	p := make(map[string]interface{})
	json.Unmarshal([]byte(data), &p)

	var chg, err = Apply(&a, p)

	assert.NoError(t, err)
	assert.True(t, chg)
	assert.Equal(t, "Darth", a.Contact.FirstName)
	assert.Equal(t, "Vader", a.Contact.LastName)
	assert.Equal(t, 100500, a.Salary)
}
