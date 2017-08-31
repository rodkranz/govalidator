package govalidator

import (
	"strconv"
	"reflect"
)

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value, nameInTag, nameAlias string, errors Errors) Errors {
	switch valueKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val == "" {
			val = "0"
		}
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			errors.Add(nameInTag, nameAlias, ERR_INTERGER_TYPE, ERR_PARSE_VALUE_PARSED_INTEGER)
		} else {
			structField.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val == "" {
			val = "0"
		}
		uintVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			errors.Add(nameInTag, nameAlias, ERR_INTERGER_TYPE, ERR_PARSE_VALUE_PARSED_UNSINED_INTEGER)
		} else {
			structField.SetUint(uintVal)
		}
	case reflect.Bool:
		if val == "on" {
			structField.SetBool(true)
			break
		}

		if val == "" {
			val = "false"
		}
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			errors.Add(nameInTag, nameAlias, ERR_BOOLEAN_TYPE, ERR_PARSE_VALUE_PARSED_BOOLEAN)
		} else if boolVal {
			structField.SetBool(true)
		}
	case reflect.Float32:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			errors.Add(nameInTag, nameAlias, ERR_FLOAT_TYPE, ERR_PARSE_VALUE_PARSED_INT_32)
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.Float64:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			errors.Add(nameInTag, nameAlias, ERR_FLOAT_TYPE, ERR_PARSE_VALUE_PARSED_INT_64)
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.String:
		structField.SetString(val)
	}
	return errors
}
