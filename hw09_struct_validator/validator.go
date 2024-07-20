package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

const vTagName = "validate"

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errMsg strings.Builder

	for _, vErr := range v {
		msg := fmt.Sprintf("field '%s': %s. ", vErr.Field, vErr.Err)
		errMsg.WriteString(msg)
	}

	return errMsg.String()
}

func Validate(v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return fmt.Errorf("v is not a struct")
	}
	var vErrors = make(ValidationErrors, 0)
	structure := reflect.ValueOf(v)
	strType := structure.Type()

	for i := 0; i < structure.NumField(); i++ {
		strFieldType := strType.Field(i)
		if strFieldType.Tag.Get(vTagName) == "" {
			continue
		}
		strFieldValue := structure.Field(i)
		if err := validate(strFieldType, strFieldValue, &vErrors); err != nil {
			return err
		}
	}

	if len(vErrors) != 0 {
		return fmt.Errorf("validation failed: %w", vErrors)
	}

	return nil
}

func validate(t reflect.StructField, v reflect.Value, vErrors *ValidationErrors) error {
	switch t.Type.Kind() {
	case reflect.Int:
		if err := validateIntField(t, int(v.Int()), vErrors); err != nil {
			return err
		}
	case reflect.String:
		if err := validateStringField(t, v.String(), vErrors); err != nil {
			return err
		}
	case reflect.Slice:
		if t.Type.Elem().Kind() == reflect.Int {
			if err := validateIntSlice(t, v, vErrors); err != nil {
				return err
			}
		} else if t.Type.Elem().Kind() == reflect.String {
			if err := validateStringSlice(t, v, vErrors); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateStringSlice(f reflect.StructField, v reflect.Value, vErrors *ValidationErrors) error {
	if strs, ok := v.Interface().([]string); ok {
		for _, s := range strs {
			if err := validateStringField(f, s, vErrors); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("type assertion to []string for field '%s' failed", f.Name)
}

func validateIntSlice(f reflect.StructField, v reflect.Value, vErrors *ValidationErrors) error {
	if ints, ok := v.Interface().([]int); ok {
		for _, i := range ints {
			if err := validateIntField(f, i, vErrors); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("type assertion to []int for field '%s' failed", f.Name)
}

func validateStringField(fInfo reflect.StructField, v string, vErrors *ValidationErrors) error {
	tags, err := parseStringTags(fInfo.Tag)
	if err != nil {
		return err
	}
	for _, t := range tags {
		if err = t.Validate(v); err != nil {
			vErr := ValidationError{
				Field: fInfo.Name,
				Err:   err,
			}
			*vErrors = append(*vErrors, vErr)
		}
	}
	return nil
}

func validateIntField(fInfo reflect.StructField, v int, vErrors *ValidationErrors) error {
	tags, err := parseIntTags(fInfo.Tag)
	if err != nil {
		return err
	}
	for _, t := range tags {
		if err = t.Validate(v); err != nil {
			vErr := ValidationError{
				Field: fInfo.Name,
				Err:   err,
			}
			*vErrors = append(*vErrors, vErr)
		}
	}
	return nil
}

func parseIntTags(t reflect.StructTag) ([]IntTag, error) {
	vTags := t.Get(vTagName)
	splitVTags := strings.Split(vTags, "|")
	tags := make([]IntTag, 0)
	for _, tag := range splitVTags {
		if intTag, err := NewIntTag(tag); err != nil {
			return nil, err
		} else {
			tags = append(tags, *intTag)
		}
	}
	return tags, nil
}

func parseStringTags(t reflect.StructTag) ([]StringTag, error) {
	vTags := t.Get(vTagName)
	splitVTags := strings.Split(vTags, "|")
	tags := make([]StringTag, 0)
	for _, tag := range splitVTags {
		if stringTag, err := NewStringTag(tag); err != nil {
			return nil, err
		} else {
			tags = append(tags, *stringTag)
		}
	}
	return tags, nil
}
