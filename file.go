package commitlint

import (
	"bufio"
	"os"

	"gopkg.in/yaml.v3"
)

// WriteConfToFile util func to write config object to given file
func WriteConfToFile(confPath string, conf *Config) (retErr error) {
	file, err := os.Create(confPath)
	if err != nil {
		return err
	}
	defer func() {
		err1 := file.Close()
		if retErr == nil && err1 != nil {
			retErr = err1
		}
	}()

	w := bufio.NewWriter(file)
	defer func() {
		err1 := w.Flush()
		if retErr == nil && err1 != nil {
			retErr = err1
		}
	}()

	enc := yaml.NewEncoder(w)
	return enc.Encode(conf)
}

// WriteHookToFile util func to write commit-msg hook to given file
func WriteHookToFile(hookFilePath string) (retErr error) {
	// if commit-msg already exists skip creating or overwriting it
	if IsFileExists(hookFilePath) {
		return nil
	}
	// commit-msg needs to be executable
	file, err := os.OpenFile(hookFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	defer func() {
		err1 := file.Close()
		if retErr == nil && err1 != nil {
			retErr = err1
		}
	}()

	commitHook := CommitMsgHook()
	_, err = file.WriteString(commitHook)
	return err
}
