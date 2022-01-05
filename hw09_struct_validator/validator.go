package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	ePack := ValidationErrors{}
	value := reflect.ValueOf(v)

	boo := value.Kind()
	_ = boo

	if value.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	t := value.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldValue := value.Field(i)
		tag, isset := t.Field(i).Tag.Lookup("validate")

		if isset {
			rules := strings.Split(tag, "|")
			ePack = append(ePack, checkRule(rules, t.Field(i).Name, fieldValue)...)
		}
	}
	return ePack
}

func checkRule(rules []string, fieldName string, fieldValue reflect.Value) ValidationErrors {
	kind := fieldValue.Kind()

	if kind == reflect.Int {
		return validateInt(rules, fieldName, fieldValue)
	}

	if kind == reflect.String {
		return validateString(rules, fieldName, fieldValue)
	}

	if kind == reflect.Slice {
		return validateSlice(rules, fieldName, fieldValue)
	}

	return nil
}

func validateSlice(rules []string, fieldName string, fieldValue reflect.Value) ValidationErrors {
	errorsPack := ValidationErrors{}

	if fieldValue.Type().String() == "[]string" {
		for i := 0; i < fieldValue.Len(); i++ {
			errorsPack = append(errorsPack, validateString(rules, fieldName, fieldValue.Index(i))...)
		}
	}

	if fieldValue.Type().String() == "[]int" {
		for i := 0; i < fieldValue.Len(); i++ {
			errorsPack = append(errorsPack, validateInt(rules, fieldName, fieldValue.Index(i))...)
		}
	}

	return errorsPack
}

func validateInt(rules []string, fieldName string, fieldValue reflect.Value) ValidationErrors {
	ePack := ValidationErrors{}

	for i := 0; i < len(rules); i++ {
		keyValue := strings.Split(rules[i], ":")
		key := keyValue[0]
		val := keyValue[1]

		switch key {
		case "min":
			vi, err := strconv.Atoi(val)
			if err != nil {
				ePack = append(ePack, ValidationError{Field: "Program error", Err: errors.New("incorrect validation rule")})
				return ePack
			}

			if fieldValue.Int() < int64(vi) {
				ePack = append(ePack, ValidationError{Field: fieldName, Err: errors.New("should be greater than " + val)})
			}
		case "max":
			vi, err := strconv.Atoi(val)
			if err != nil {
				ePack = append(ePack, ValidationError{Field: "Program error", Err: errors.New("incorrect validation rule")})
				return ePack
			}

			if fieldValue.Int() > int64(vi) {
				ePack = append(ePack, ValidationError{Field: fieldName, Err: errors.New("should be less than " + val)})
			}
		case "in":
			has := false
			elems := strings.Split(val, ",")
			for i := 0; i < len(elems); i++ {
				vi, err := strconv.Atoi(elems[i])
				if err != nil {
					ePack = append(ePack, ValidationError{Field: "Program error", Err: errors.New("incorrect validation rule")})
					return ePack
				}

				if fieldValue.Int() == int64(vi) {
					has = true
					break
				}
			}
			if !has {
				ePack = append(ePack, ValidationError{Field: fieldName, Err: errors.New("should be in " + val)})
			}
		default:
			ePack = append(ePack, ValidationError{Field: "Program error", Err: errors.New("validation rule not found")})
			return ePack
		}
	}

	return ePack
}

func validateString(rules []string, fieldName string, fieldValue reflect.Value) ValidationErrors {
	ePack := ValidationErrors{}

	for i := 0; i < len(rules); i++ {
		keyValue := strings.Split(rules[i], ":")
		key := keyValue[0]
		val := keyValue[1]

		switch key {
		case "len":
			vl, err := strconv.Atoi(val)
			if err != nil {
				ePack = append(ePack, ValidationError{Field: "Program error", Err: errors.New("incorrect validation rule")})
				return ePack
			}

			if len(fieldValue.String()) != vl {
				ePack = append(ePack, ValidationError{Field: fieldName, Err: errors.New("len should be " + val)})
			}
		case "in":
			in := false
			elems := strings.Split(val, ",")
			for i := 0; i < len(elems); i++ {
				if fieldValue.String() == elems[i] {
					in = true
					break
				}
			}
			if !in {
				ePack = append(ePack, ValidationError{Field: fieldName, Err: errors.New("should be in " + val)})
			}
		case "regexp":
			r, _ := regexp.Compile(val)
			result := r.MatchString(fieldValue.String())

			if !result {
				ePack = append(ePack, ValidationError{Field: fieldName, Err: errors.New("should fit " + val)})
			}
		default:
			ePack = append(ePack, ValidationError{Field: "Program error", Err: errors.New("validation rule not found")})
			return ePack
		}
	}

	return ePack
}
