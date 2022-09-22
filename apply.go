package patch

import (
	"fmt"
	"github.com/fatih/structs"
	"reflect"
	"strings"
)

// Apply updates the target struct in-place with non-zero values from the patch struct.
// Only fields with the same name and type get updated. Fields in the patch struct can be
// pointers to the target's type.
//
// Returns true if any value has been changed.
func Apply(target interface{}, patch map[string]interface{}) (changed bool, err error) {
	var dst = structs.New(target)

	for key, value := range patch {
		var name = key
		var dstField, ok = findField(dst, name)
		if !ok {
			continue // skip non-existing fields
		}

		if dstField.Kind() == reflect.Struct ||
			(dstField.Kind() == reflect.Pointer &&
				reflect.Indirect(reflect.ValueOf(dstField.Value())).Kind() == reflect.Struct) {
			// recursive for structs and pointers to structs
			iChanged, iErr := Apply(dstField.Value(), value.(map[string]interface{}))
			if iErr != nil {
				err = iErr
				return
			}
			changed = changed || iChanged
			continue
		}

		var srcValue = reflect.ValueOf(value)
		srcValue = reflect.Indirect(srcValue)
		if skind, dkind := srcValue.Kind(), dstField.Kind(); skind != dkind {
			srcValue, err = casting(value, dstField)
			if err != nil {
				return
			}
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

func casting(src interface{}, dst *structs.Field) (reflect.Value, error) {
	switch dst.Kind() {
	case reflect.Int:
		f, ok := src.(float64)
		if ok {
			i := int(f)
			return reflect.Indirect(reflect.ValueOf(i)), nil
		}
	case reflect.Int8:
		f, ok := src.(float64)
		if ok {
			i := int8(f)
			return reflect.Indirect(reflect.ValueOf(i)), nil
		}
	case reflect.Int16:
		f, ok := src.(float64)
		if ok {
			i := int16(f)
			return reflect.Indirect(reflect.ValueOf(i)), nil
		}
	case reflect.Int32:
		f, ok := src.(float64)
		if ok {
			i := int32(f)
			return reflect.Indirect(reflect.ValueOf(i)), nil
		}
	case reflect.Int64:
		f, ok := src.(float64)
		if ok {
			i := int64(f)
			return reflect.Indirect(reflect.ValueOf(i)), nil
		}
	}
	return reflect.Value{}, fmt.Errorf("field `%v` types mismatch while patching: %v vs %v", src, dst.Kind(), reflect.ValueOf(src).Kind())
}

func findField(dst *structs.Struct, name string) (*structs.Field, bool) {
	for _, field := range dst.Fields() {
		tag := field.Tag("json")
		if tag == "" {
			tag = field.Name()
		} else {
			tag, _, _ = strings.Cut(tag, ",")
		}
		if tag == name {
			if field.IsExported() {
				return field, true
			}
		}
	}
	return nil, false
}
