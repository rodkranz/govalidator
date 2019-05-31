// Package go_validator
package govalidator

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Validate(obj interface{}) (bool, Errors) {
	var errs Errors
	v := reflect.ValueOf(obj)
	k := v.Kind()
	if k == reflect.Interface || k == reflect.Ptr {
		v = v.Elem()
		k = v.Kind()
	}

	if k == reflect.Slice || k == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			e := v.Index(i).Interface()
			errs = validateStruct(errs, e)
			if validator, ok := e.(Validator); ok {
				errs = validator.Validate(errs)
			}
		}
	} else {
		errs = validateStruct(errs, obj)
		if validator, ok := obj.(Validator); ok {
			errs = validator.Validate(errs)
		}
	}

	return len(errs) == 0, errs
}

// Performs required field checking on a struct
func validateStruct(errors Errors, obj interface{}) Errors {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return errors
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.Tag.Get("json") == "-" || !val.Field(i).CanInterface() {
			continue
		}

		fieldVal := val.Field(i)
		fieldValue := fieldVal.Interface()
		zero := reflect.Zero(field.Type).Interface()

		if field.Type.Kind() == reflect.Struct ||
			(field.Type.Kind() == reflect.Ptr && !reflect.DeepEqual(zero, fieldValue) &&
				field.Type.Elem().Kind() == reflect.Struct) {
			errors = validateStruct(errors, fieldValue)
		}
		errors = validateField(errors, zero, field, fieldVal, fieldValue)
	}
	return errors
}

func validateField(errors Errors, zero interface{}, field reflect.StructField, fieldVal reflect.Value, fieldValue interface{}) Errors {
	if fieldVal.Kind() == reflect.Slice {
		for i := 0; i < fieldVal.Len(); i++ {
			sliceVal := fieldVal.Index(i)
			if sliceVal.Kind() == reflect.Ptr {
				sliceVal = sliceVal.Elem()
			}

			sliceValue := sliceVal.Interface()
			zero := reflect.Zero(sliceVal.Type()).Interface()
			if sliceVal.Kind() == reflect.Struct ||
				(sliceVal.Kind() == reflect.Ptr && !reflect.DeepEqual(zero, sliceValue) &&
					sliceVal.Elem().Kind() == reflect.Struct) {
				errors = validateStruct(errors, sliceValue)
			}
		}
	}

VALIDATE_RULES:
	for _, rule := range strings.Split(field.Tag.Get("validate"), ";") {
		if len(rule) == 0 {
			continue
		}

		alias := field.Tag.Get("alias")

		switch {
		case rule == OMIT_EMPTY:
			if reflect.DeepEqual(zero, fieldValue) {
				break VALIDATE_RULES
			}

		case rule == REQUIRED:
			v := reflect.ValueOf(fieldValue)
			if v.Kind() == reflect.Slice {
				if v.Len() == 0 {
					errors.Add(field.Name, alias, ERR_REQUIRED, REQUIRED)
					break VALIDATE_RULES
				}

				continue
			}

			if reflect.DeepEqual(zero, fieldValue) {
				errors.Add(field.Name, alias, ERR_REQUIRED, REQUIRED)
				break VALIDATE_RULES
			}

		case rule == ALPHA_DASH:
			if AlphaDashPattern.MatchString(fmt.Sprintf("%v", fieldValue)) {
				errors.Add(field.Name, alias, ERR_ALPHA_DASH, ALPHA_DASH)
				break VALIDATE_RULES
			}

		case rule == ALPHA_DASH_DOT:
			if AlphaDashDotPattern.MatchString(fmt.Sprintf("%v", fieldValue)) {
				errors.Add(field.Name, alias, ERR_ALPHA_DASH_DOT, ALPHA_DASH_DOT)
				break VALIDATE_RULES
			}

		case rule == NUMERIC:
			if NumericPattern.MatchString(fmt.Sprintf("%v", fieldValue)) {
				errors.Add(field.Name, alias, ERR_NUMERIC, NUMERIC)
				break VALIDATE_RULES
			}

		case rule == EMAIL:
			if !EmailPattern.MatchString(fmt.Sprintf("%v", fieldValue)) {
				errors.Add(field.Name, alias, ERR_EMAIL, EMAIL)
				break VALIDATE_RULES
			}

		case rule == URL:
			str := fmt.Sprintf("%v", fieldValue)
			if len(str) == 0 {
				continue
			} else if !isURL(str) {
				errors.Add(field.Name, alias, ERR_URL, URL)
				break VALIDATE_RULES
			}

		case strings.HasPrefix(rule, SIZE+"("):
			size, _ := strconv.Atoi(rule[5 : len(rule)-1])
			if str, ok := fieldValue.(string); ok && utf8.RuneCountInString(str) != size {
				errors.Add(field.Name, alias, ERR_SIZE, SIZE)
				break VALIDATE_RULES
			}
			v := reflect.ValueOf(fieldValue)
			if v.Kind() == reflect.Slice && v.Len() != size {
				errors.Add(field.Name, alias, ERR_SIZE, SIZE)
				break VALIDATE_RULES
			}

			if v.Kind() == reflect.Ptr {
				elm := v.Elem()
				if elm.Kind() == reflect.String && elm.Len() != size {
					errors.Add(field.Name, alias, ERR_MAX_SIZE, MAX_SIZE)
					break VALIDATE_RULES
				}
			}
		case strings.HasPrefix(rule, MIN_SIZE+"("):
			min, _ := strconv.Atoi(rule[8 : len(rule)-1])
			if str, ok := fieldValue.(string); ok && utf8.RuneCountInString(str) < min {
				errors.Add(field.Name, alias, ERR_MIN_SIZE, MIN_SIZE)
				break VALIDATE_RULES
			}
			v := reflect.ValueOf(fieldValue)
			if v.Kind() == reflect.Slice && v.Len() < min {
				errors.Add(field.Name, alias, ERR_MIN_SIZE, MIN_SIZE)
				break VALIDATE_RULES
			}

			if v.Kind() == reflect.Ptr {
				elm := v.Elem()
				if elm.Kind() == reflect.String && elm.Len() < min {
					errors.Add(field.Name, alias, ERR_MIN_SIZE, MIN_SIZE)
					break VALIDATE_RULES
				}
			}
		case strings.HasPrefix(rule, MAX_SIZE+"("):
			max, _ := strconv.Atoi(rule[8 : len(rule)-1])
			if str, ok := fieldValue.(string); ok && utf8.RuneCountInString(str) > max {
				errors.Add(field.Name, alias, ERR_MAX_SIZE, MAX_SIZE)
				break VALIDATE_RULES
			}
			v := reflect.ValueOf(fieldValue)
			if v.Kind() == reflect.Slice && v.Len() > max {
				errors.Add(field.Name, alias, ERR_MAX_SIZE, MAX_SIZE)
				break VALIDATE_RULES
			}

			if v.Kind() == reflect.Ptr {
				elm := v.Elem()
				if elm.Kind() == reflect.String && elm.Len() > max {
					errors.Add(field.Name, alias, ERR_MAX_SIZE, MAX_SIZE)
					break VALIDATE_RULES
				}
			}
		case strings.HasPrefix(rule, INCLUDE+"("):
			if !strings.Contains(fmt.Sprintf("%v", fieldValue), rule[8:len(rule)-1]) {
				errors.Add(field.Name, alias, ERR_INCLUDE, INCLUDE)
				break VALIDATE_RULES
			}

		case strings.HasPrefix(rule, EXCLUDE+"("):
			if strings.Contains(fmt.Sprintf("%v", fieldValue), rule[8:len(rule)-1]) {
				errors.Add(field.Name, alias, ERR_EXCLUDE, EXCLUDE)
				break VALIDATE_RULES
			}
		case strings.HasPrefix(rule, DEFAULT+"("):
			if reflect.DeepEqual(zero, fieldValue) {
				if fieldVal.CanAddr() {
					errors = setWithProperType(field.Type.Kind(), rule[8:len(rule)-1], fieldVal, field.Tag.Get("json"), alias, errors)
				} else {
					errors.Add(field.Name, alias, ERR_EXCLUDE, DEFAULT)
					break VALIDATE_RULES
				}
			}
		default:
			// Apply custom validation rules
			var isValid bool
			for i := range ruleMapper {
				if ruleMapper[i].IsMatch(rule) {
					isValid, errors = ruleMapper[i].IsValid(errors, field.Name, alias, fieldValue)
					if !isValid {
						break VALIDATE_RULES
					}
				}
			}
			for i := range paramRuleMapper {
				if paramRuleMapper[i].IsMatch(rule) {
					isValid, errors = paramRuleMapper[i].IsValid(errors, rule, field.Name, alias, fieldValue)
					if !isValid {
						break VALIDATE_RULES
					}
				}
			}
		}
	}
	return errors
}

// IsURL check if the string is an URL.
func isURL(str string) bool {
	if str == "" || utf8.RuneCountInString(str) >= MAX_URL_RUNE_COUNT || len(str) <= MIN_URL_RUNE_COUNT || strings.HasPrefix(str, ".") {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	return URLPattern.MatchString(str)

}
