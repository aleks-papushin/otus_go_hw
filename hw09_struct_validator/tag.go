package hw09structvalidator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type intCheckType string

type stringCheckType string

const (
	minimum intCheckType = "min"
	maximum intCheckType = "max"
	intSet  intCheckType = "in"

	length stringCheckType = "length"
	regExp stringCheckType = "regExp"
	strSet stringCheckType = "in"
)

type IntTag struct {
	checkType intCheckType
	values    []string
}

func NewIntTag(s string) (*IntTag, error) {
	tokens := strings.Split(s, ":")
	if len(tokens) != 2 {
		err := fmt.Errorf("wrong validation tag format for tag %s\n", s)
		return nil, err
	}
	checkTypeName := tokens[0]
	values := strings.Split(tokens[1], ",")
	if !isValidIntCheckType(checkTypeName) {
		err := fmt.Errorf("wrong int check type (%s). Valid are: %s, %s, %s", checkTypeName, minimum, maximum, intSet)
		return nil, err
	}
	checkType := intCheckType(checkTypeName)
	if !areValidFor(checkType, values) {
		err := fmt.Errorf("wrong value(-s) for int check type %s (%s)", checkType, values)
		return nil, err
	}

	return &IntTag{
		checkType: checkType,
		values:    values,
	}, nil
}

func (t *IntTag) Validate(v int) error {
	switch t.checkType {
	case minimum:
		if minTagValue, _ := strconv.Atoi(t.values[0]); v < minTagValue {
			return fmt.Errorf("value (%d) less than minmum (%d)", v, minTagValue)
		}
		return nil
	case maximum:
		if maxTagValue, _ := strconv.Atoi(t.values[0]); v > maxTagValue {
			return fmt.Errorf("value (%d) more than maximum (%d)", v, maxTagValue)
		}
		return nil
	case intSet:
		for _, intSetValue := range t.values {
			if i, _ := strconv.Atoi(intSetValue); v == i {
				return nil
			}
		}
		return fmt.Errorf("value (%d) is not in set (%s)", v, t.values)
	default:
		return fmt.Errorf("unexpected checkType %s", t.checkType)
	}
}

func areValidFor(checkType intCheckType, values []string) bool {
	switch checkType {
	case minimum, maximum:
		isSingleElement := len(values) == 1
		_, err := strconv.Atoi(values[0])
		return isSingleElement && err == nil
	case intSet:
		for _, v := range values {
			if _, err := strconv.Atoi(v); err != nil {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func isValidIntCheckType(s string) bool {
	switch intCheckType(s) {
	case minimum, maximum, intSet:
		return true
	default:
		return false
	}
}

type StringTag struct {
	checkType stringCheckType
	values    []string
}

func NewStringTag(s string) (*StringTag, error) {
	tokens := strings.Split(s, ":")
	if len(tokens) != 2 {
		err := fmt.Errorf("wrong validation tag format for tag %s\n", s)
		return nil, err
	}
	checkTypeName := tokens[0]
	if !isValidStringCheckType(checkTypeName) {
		err := fmt.Errorf("wrong string check type (%s). Valid are: %s, %s, %s", checkTypeName, length, regExp, strSet)
		return nil, err
	}
	checkType := stringCheckType(checkTypeName)
	values := make([]string, 0)
	if checkType == strSet {
		values = strings.Split(tokens[1], ",")
	} else if checkType == regExp {
		re := tokens[1]
		_, err := regexp.Compile(re)
		if err != nil {
			return nil, fmt.Errorf("invalid regexp provided: %s", re)
		}
		values = append(values, re)
	} else {
		values = append(values, tokens[1])
	}
	if !areValidForString(checkType, values) {
		err := fmt.Errorf("wrong value(-s) for string check type %s (%s)", checkType, values)
		return nil, err
	}

	return &StringTag{
		checkType: checkType,
		values:    values,
	}, nil
}

func (t *StringTag) Validate(s string) error {
	switch t.checkType {
	case length:
		l, _ := strconv.Atoi(t.values[0])
		if len(s) != l {
			return fmt.Errorf("length is not equal to %d", l)
		}
		return nil
	case regExp:
		re := t.values[0]
		regExp, _ := regexp.Compile(re) // regexp validation already passed in tag constructor
		if ok := regExp.Match([]byte(s)); !ok {
			return fmt.Errorf("value (%s) is not matched with regexp %s", s, re)
		}
		return nil
	case strSet:
		for _, v := range t.values {
			if s == v {
				return nil
			}
		}
		return fmt.Errorf("value (%s) is not in set (%s)", s, t.values)
	default:
		return fmt.Errorf("unexpected checkType %s", t.checkType)
	}
}

func isValidStringCheckType(s string) bool {
	switch stringCheckType(s) {
	case length, regExp, strSet:
		return true
	default:
		return false
	}
}

func areValidForString(checkType stringCheckType, values []string) bool {
	isSingleElement := len(values) == 1
	switch checkType {
	case length:
		_, err := strconv.Atoi(values[0])
		return isSingleElement && err == nil
	case regExp, strSet:
		return true
	default:
		return false
	}
}
