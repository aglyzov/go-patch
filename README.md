# go-patch

With [go-patch](https://github.com/aglyzov/go-patch) you can selectively update [golang](http://golang.org) `structs` with
values from other structs.

### Example
```go
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
    NonExisting []byte  // fields not present in the target are also ignored
}

var a = Target{
    FirstName: "Anakin",
    LastName:  "Skywalker",
    Salary:    123,
    Extra:     "unchanged",
}
var p = Patch{
    FirstName: "Darth",
    LastName:  "Vader",
}
var changed, err = patch.Struct(&a, p)

// now `a` is {"Darth", "Vader", 123, "unchanged"}
```
