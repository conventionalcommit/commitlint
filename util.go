package commitlint

import "os"

// IsFileExists checks if given file exists
func IsFileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}
