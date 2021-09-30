package rule

import (
	"fmt"
)

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
