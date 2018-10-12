# go-patch

With [go-patch](https://github.com/aglyzov/go-patch) you can selectively update [golang](http://golang.org) `structs` with
values from other structs.

## API

### func Struct(dstStruct, srcStruct interface{}) (bool, error)
`patch.Struct` updates a destination structure in-place with the same name fields
from a supplied patch struct. Fields get matched by their names (case sensitive).
Thus patch fields with a name not present in the destination structure are ignored. 
`Zero-value` and `unexported` fields in a patch struct are also ignored.

Notice, both the destination and patch structure can have `embedded` structs in them.


### Example
```go
import "github.com/aglyzov/go-patch"

type Employee struct {
    FirstName string
    LastName  string
    Salary    int
    Extra     string
}
type Patch struct {
    FirstName  string   // names and types should match
    LastName   *string  // however a patching field can also be a pointer 
    Salary     int      // only non-zero values are considered

    unexported bool     // unexported fields are ignored 
    Unknown    []byte   // fields not present in the target are also ignored
}

var e = Employee{
    FirstName: "Anakin",
    LastName:  "Skywalker",
    Salary:    123,
    Extra:     "unchanged",
}
var lastName = "Vader"
var p = Patch{
    FirstName:  "Darth",
    LastName:   &lastName, // pointer to a string
    Salary:     0,         // zero-value is ignored
    unexported: true,
    Unknown:    []byte("ignored"),
}
var changed, err = patch.Struct(&e, p)

// now `e` is {"Darth", "Vader", 123, "unchanged"}
```
