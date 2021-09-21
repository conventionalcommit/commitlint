package rule

import (
	"fmt"
)

func setBoolArg(retVal *bool, arg interface{}, ruleName string) error {
	boolVal, err := toBool(arg, ruleName)
	if err != nil {
		return err
	}
	*retVal = boolVal
	return nil
}

func setIntArg(retVal *int, arg interface{}, ruleName string) error {
	intVal, err := toInt(arg, ruleName)
	if err != nil {
		return err
	}
	*retVal = intVal
	return nil
}

func setStringArg(retVal *string, arg interface{}, ruleName string) error {
	strVal, err := toString(arg, ruleName)
	if err != nil {
		return err
	}
	*retVal = strVal
	return nil
}

func setStringArrArg(retVal *[]string, arg interface{}, ruleName string) error {
	arrVal, err := toStringArr(arg, ruleName)
	if err != nil {
		return err
	}
	*retVal = arrVal
	return nil
}

func toBool(arg interface{}, ruleName string) (bool, error) {
	boolVal, ok := arg.(bool)
	if !ok {
		return false, fmt.Errorf("%s expects bool value, but got %#v", ruleName, arg)
	}
	return boolVal, nil
}

func toInt(arg interface{}, ruleName string) (int, error) {
	intVal, ok := arg.(int)
	if !ok {
		return 0, fmt.Errorf("%s expects int value, but got %#v", ruleName, arg)
	}
	return intVal, nil
}

func toString(arg interface{}, ruleName string) (string, error) {
	strVal, ok := arg.(string)
	if !ok {
		return "", fmt.Errorf("%s expects string value, but got %#v", ruleName, arg)
	}
	return strVal, nil
}

func toStringArr(arg interface{}, ruleName string) ([]string, error) {
	switch argVal := arg.(type) {
	case []interface{}:
		strArr := make([]string, len(argVal))
		for index, a := range argVal {
			strVal, ok := a.(string)
			if !ok {
				return nil, fmt.Errorf("%s expects array of string value, but got %#v", ruleName, arg)
			}
			strArr[index] = strVal
		}
		return strArr, nil
	case []string:
		return argVal, nil
	}
	return nil, fmt.Errorf("%s expects array of string value, but got %#v", ruleName, arg)
}
