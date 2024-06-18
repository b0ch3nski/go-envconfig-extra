package envconfigext

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/wkhere/dtf"
)

// StructFieldScan is an opinionated struct to string converter with optional field value masking made especially for
// pretty printing configurations. It's not performance optimized so should be used with care.
func StructFieldScan(iface any) string {
	fields := make([]string, 0)
	ifVal := reflect.ValueOf(iface)
	ifType := reflect.TypeOf(iface)

	for i := 0; i < ifType.NumField(); i++ {
		value := ifVal.Field(i)
		if !value.CanInterface() {
			continue
		}
		field := ifType.Field(i)
		str := anyTypeToStr(value)

		if tagVal, ok := field.Tag.Lookup("secret"); ok {
			maskType, maskParam, _ := strings.Cut(tagVal, "=")
			switch maskType {
			case "", "redact":
				str = "REDACTED"
			case "mask":
				leaveChars, _ := strconv.ParseUint(maskParam, 10, 0)
				str = maskLeft(str, leaveChars)
			}
		}
		fields = append(fields, fmt.Sprintf("%s=%s", field.Name, str))
	}
	return strings.Join(fields, " | ")
}

func anyTypeToStr(val reflect.Value) string {
	switch i := val.Interface().(type) {
	case time.Time:
		return i.Format(time.RFC3339)
	case time.Duration:
		return dtf.Fmt(i)
	case fmt.Stringer:
		return fmt.Sprintf("{%s}", i.String())
	}

	switch val.Kind() {
	case reflect.Pointer, reflect.Interface:
		if val.IsNil() {
			return "nil"
		}
		return anyTypeToStr(val.Elem())
	case reflect.Struct:
		return fmt.Sprintf("{%s}", StructFieldScan(val.Interface()))
	case reflect.Map:
		flatten := make([]string, 0, val.Len())
		iter := val.MapRange()
		for iter.Next() {
			k, v := anyTypeToStr(iter.Key()), anyTypeToStr(iter.Value())
			flatten = append(flatten, fmt.Sprintf("%s=%s", k, v))
		}
		return fmt.Sprintf("[%s]", strings.Join(flatten, ", "))
	case reflect.Array, reflect.Slice:
		if val.Type().ConvertibleTo(reflect.TypeOf([]byte(nil))) && val.Len() > 0 {
			return "<binary data>"
		}
		flatten := make([]string, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			flatten = append(flatten, anyTypeToStr(val.Index(i)))
		}
		return fmt.Sprintf("[%s]", strings.Join(flatten, ", "))
	}

	return fmt.Sprintf("%v", val)
}

// idea based on: https://stackoverflow.com/a/48812623
func maskLeft(str string, leaveChars uint64) string {
	rs := []rune(str)
	for i := 0; i < len(rs)-int(leaveChars); i++ {
		rs[i] = 'X'
	}
	return string(rs)
}
