package remember

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"reflect"
	"runtime"
	"strings"
)

// CreateKey will generate a key based on the input arguments.
// When prefix is true, the caller's name will be used to prefix the key in an attempt to make it unique.
// The args can also be separated using sep. visual performs no functionality. It is used at code level
// to visually see how the key is structured.
func CreateKey(prefix bool, sep string, visual string, args ...interface{}) string {
	var output string

	if prefix {
		pc, file, line, ok := runtime.Caller(1)
		if !ok {
			return fmt.Sprint(args...)
		}
		details := runtime.FuncForPC(pc)
		output = fmt.Sprintf("%s_%s_%d_", details.Name(), file, line)
	}

	if sep == "" {
		output = output + fmt.Sprint(args...)
	} else {
		for i, v := range args {
			if i != 0 {
				output = output + sep
			}
			output = output + fmt.Sprint(v)
		}
	}

	return output
}

// Hash returns a crc32 hashed version of key.
func Hash(key string) string {
	return fmt.Sprintf("%08x\n", crc32.ChecksumIEEE([]byte(key)))
}

// CreateKeyStruct generates a key by converting a struct into a JSON object.
func CreateKeyStruct(strct interface{}) string {
	out := map[string]interface{}{}

	// Encode nil immediately
	if strct == nil {
		return ""
	}

	s := reflect.ValueOf(strct)

	// Check if s is a pointer
	if s.Kind() == reflect.Ptr {
		s = reflect.Indirect(s)
	}
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := typeOfT.Field(i)

		if f.PkgPath != "" {
			// Not exported
			continue
		}

		fieldName := typeOfT.Field(i).Name
		fieldTag := f.Tag.Get("json")
		fieldValRaw := s.Field(i)
		fieldVal := fieldValRaw.Interface()

		// Ignore slices
		if fieldValRaw.Kind() == reflect.Slice {
			continue
		}

		// Check if json parser would ordinarily hide the value anyway
		if fieldTag == "-" || (strings.HasSuffix(fieldTag, ",omitempty") && reflect.DeepEqual(fieldVal, reflect.Zero(reflect.TypeOf(fieldVal)).Interface())) {
			continue
		}

		if fieldTag == "" {
			out[fieldName] = fieldVal
		} else {
			out[strings.TrimSuffix(fieldTag, ",omitempty")] = fieldVal
		}
	}

	b, _ := json.Marshal(out)

	return string(b)
}
