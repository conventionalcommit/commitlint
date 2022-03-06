// Package rule contains lint rules
package rule

import (
	"fmt"
	"sort"
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

func errInvalidArg(ruleName string, err error) error {
	return fmt.Errorf("%s: invalid argument: %v", ruleName, err)
}

func errInvalidFlag(ruleName, flagName string, err error) error {
	return fmt.Errorf("%s: invalid flag '%s': %v", ruleName, flagName, err)
}

func formMinLenMsg(typ string, actualLen, expectedLen int) string {
	return fmt.Sprintf("%s length is %d, should have atleast %d chars", typ, actualLen, expectedLen)
}

func formMaxLenDesc(typ string, actualLen, expectedLen int) string {
	return fmt.Sprintf("%s length is %d, should have less than %d chars", typ, actualLen, expectedLen)
}

func formMaxLineLenDesc(typ string, expectedLen int) string {
	return fmt.Sprintf("each %s line should have less than %d chars", typ, expectedLen)
}

func search(arr []string, toFind string) bool {
	ind := sort.Search(len(arr), func(i int) bool {
		return arr[i] >= toFind
	})
	return ind < len(arr) && arr[ind] == toFind
}

func validateCharset(allowedCharset, toCheck string) (string, bool) {
	invalidRunes := ""
	for _, ch := range toCheck {
		if !strings.ContainsRune(allowedCharset, ch) {
			invalidRunes += string(ch)
		}
	}

	if invalidRunes == "" {
		return "", true
	}
	return invalidRunes, false
}

func validateMinLen(typ string, expectedLen int, toCheck string) (*lint.Issue, bool) {
	actualLen := len(toCheck)
	if actualLen >= expectedLen {
		return nil, true
	}

	desc := formMinLenMsg(typ, actualLen, expectedLen)
	return lint.NewIssue(desc), false
}

func validateMaxLen(typ string, expectedLen int, toCheck string) (*lint.Issue, bool) {
	if expectedLen < 0 {
		return nil, true
	}

	if len(toCheck) <= expectedLen {
		return nil, true
	}

	desc := formMaxLenDesc(typ, len(toCheck), expectedLen)
	return lint.NewIssue(desc), false
}

func validateMaxLineLength(typ string, expectedLen int, toCheck string) (*lint.Issue, bool) {
	lines := strings.Split(toCheck, "\n")

	msgs := []string{}
	for index, line := range lines {
		actualLen := len(line)
		if actualLen > expectedLen {
			errMsg := fmt.Sprintf("in line %d, length is %d", index+1, actualLen)
			msgs = append(msgs, errMsg)
		}
	}

	if len(msgs) == 0 {
		return nil, true
	}

	desc := formMaxLineLenDesc(typ, expectedLen)
	return lint.NewIssue(desc, msgs...), false
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
	default:
		return nil, fmt.Errorf("expects array of string value, but got %#v", arg)
	}
}
