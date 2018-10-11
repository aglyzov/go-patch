package patch

import (
    "fmt"
    "reflect"

	"github.com/fatih/structs"
)


// Struct updates the target struct in-place with non-zero values from the patch struct.
// Only fields with the same name and type get updated. Fields in the patch struct can be
// pointers to the target's type.
//
// Returns true if any value has been changed.
func Struct(target, patch interface{}) (changed bool, err error) {

    var dst = structs.New(target)

    for _, srcField := range structs.New(patch).Fields() {
        if ! srcField.IsExported() || srcField.IsZero() {
            continue  // skip unexported and zero-value fields
        }
        var name = srcField.Name()

        var dstField, ok = dst.FieldOk(name)
        if !ok {
            continue  // skip non-existing fields
        }
        var srcValue = reflect.ValueOf(srcField.Value())
        srcValue = reflect.Indirect(srcValue)
        if skind, dkind := srcValue.Kind(), dstField.Kind(); skind != dkind {
            err = fmt.Errorf("field `%v` types mismatch while patching: %v vs %v", name, dkind, skind)
            return
        }

        if ! reflect.DeepEqual(srcValue.Interface(), dstField.Value()) {
            changed = true
        }

        err = dstField.Set(srcValue.Interface())
        if err != nil {
            return
        }
    }

    return
}
