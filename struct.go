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
	var fields = structs.New(patch).Fields() // work stack

	// assuming same order for both fields
	for N := len(fields); N > 0; N = len(fields) {
		var srcField = fields[N-1] // pop the top
		fields = fields[:N-1]

		if !srcField.IsExported() {
			continue // skip unexported fields
		}
		if srcField.IsEmbedded() {
			// add the embedded fields into the work stack
			fields = append(fields, srcField.Fields()...)
			continue
		}
		if srcField.IsZero() {
			continue // skip zero-value fields
		}
		var name = srcField.Name()
		var dstField, ok = dst.FieldOk(name)
		if !ok {
			continue // skip non-existing fields
		}

		if dstField.Kind() == reflect.Struct ||
			(dstField.Kind() == reflect.Pointer &&
				reflect.Indirect(reflect.ValueOf(dstField.Value())).Kind() == reflect.Struct) {
			// recursive for structs and pointers to structs
			iChanged, iErr := Struct(dstField.Value(), srcField.Value())
			if iErr != nil {
				err = iErr
				return
			}
			changed = changed || iChanged
			continue
		}

		var srcValue = reflect.ValueOf(srcField.Value())
		srcValue = reflect.Indirect(srcValue)
		if skind, dkind := srcValue.Kind(), dstField.Kind(); skind != dkind {
			err = fmt.Errorf("field `%v` types mismatch while patching: %v vs %v", name, dkind, skind)
			return
		}

		if !reflect.DeepEqual(srcValue.Interface(), dstField.Value()) {
			changed = true
		}

		err = dstField.Set(srcValue.Interface())
		if err != nil {
			return
		}
	}
	return
}
