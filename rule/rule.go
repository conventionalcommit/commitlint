// Package rule contains lint rules
package rule

import (
	"fmt"
	"sort"
	"strings"
)

func errInvalidArg(ruleName string, err error) error {
	return fmt.Errorf("%s: invalid argument: %v", ruleName, err)
}

func errInvalidFlag(ruleName, flagName string, err error) error {
	return fmt.Errorf("%s: invalid flag '%s': %v", ruleName, flagName, err)
}

func checkCharset(charset, toCheck string) (string, bool) {
	invalidRunes := ""
	for _, ch := range toCheck {
		if !strings.ContainsRune(charset, ch) {
			invalidRunes += string(ch)
		}
	}

	if invalidRunes == "" {
		return "", true
	}

	return invalidRunes, false
}

func search(arr []string, toFind string) bool {
	ind := sort.Search(len(arr), func(i int) bool {
		return arr[i] >= toFind
	})
	return ind < len(arr) && arr[ind] == toFind
}

func checkMaxLen(checkLen int, toCheck string) ([]string, bool) {
	if checkLen < 0 {
		return nil, true
	}

	if len(toCheck) <= checkLen {
		return nil, true
	}

	errMsg := fmt.Sprintf("length is %d, should have less than %d chars", len(toCheck), checkLen)
	return []string{errMsg}, false
}

func checkMaxLineLength(checkLen int, toCheck string) ([]string, bool) {
	lines := strings.Split(toCheck, "\n")

	msgs := []string{}
	for index, line := range lines {
		actualLen := len(line)
		if actualLen > checkLen {
			errMsg := fmt.Sprintf("in line %d, length is %d, should have less than %d chars", index+1, actualLen, checkLen)
			msgs = append(msgs, errMsg)
		}
	}

	if len(msgs) == 0 {
		return nil, true
	}
	return msgs, false
}

func checkMinLen(checkLen int, toCheck string) ([]string, bool) {
	actualLen := len(toCheck)
	if actualLen >= checkLen {
		return nil, true
	}
	errMsg := fmt.Sprintf("length is %d, should have atleast %d chars", actualLen, checkLen)
	return []string{errMsg}, false
}

func setBoolArg(retVal *bool, arg interface{}) error {
	boolVal, err := toBool(arg)
	if err != nil {
		return err
	}
	*retVal = boolVal
	return nil
}

func setIntArg(retVal *int, arg interface{}) error {
	intVal, err := toInt(arg)
	if err != nil {
		return err
	}
	*retVal = intVal
	return nil
}

func setStringArg(retVal *string, arg interface{}) error {
	strVal, err := toString(arg)
	if err != nil {
		return err
	}
	*retVal = strVal
	return nil
}

func setStringArrArg(retVal *[]string, arg interface{}) error {
	arrVal, err := toStringArr(arg)
	if err != nil {
		return err
	}
	*retVal = arrVal
	return nil
}

func toBool(arg interface{}) (bool, error) {
	boolVal, ok := arg.(bool)
	if !ok {
		return false, fmt.Errorf("expects bool value, but got %#v", arg)
	}
	return boolVal, nil
}

func toInt(arg interface{}) (int, error) {
	intVal, ok := arg.(int)
	if !ok {
		return 0, fmt.Errorf("expects int value, but got %#v", arg)
	}
	return intVal, nil
}

func toString(arg interface{}) (string, error) {
	strVal, ok := arg.(string)
	if !ok {
		return "", fmt.Errorf("expects string value, but got %#v", arg)
	}
	return strVal, nil
}

func toStringArr(arg interface{}) ([]string, error) {
	switch argVal := arg.(type) {
	case []interface{}:
		strArr := make([]string, len(argVal))
		for index, a := range argVal {
			strVal, ok := a.(string)
			if !ok {
				return nil, fmt.Errorf("expects array of string value, but got %#v", arg)
			}
			strArr[index] = strVal
		}
		return strArr, nil
	case []string:
		return argVal, nil
	}
	return nil, fmt.Errorf("expects array of string value, but got %#v", arg)
}
