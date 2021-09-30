package rule

import "fmt"

func errInvalidArg(ruleName string, err error) error {
	return fmt.Errorf("%s: invalid argument: %v", ruleName, err)
}

func errInvalidFlag(ruleName, flagName string, err error) error {
	return fmt.Errorf("%s: invalid flag '%s': %v", ruleName, flagName, err)
}
