package internal

import (
	"errors"
	"os"
	"path/filepath"
)

const CommitlintConfigEnv = "COMMITLINT_CONFIG"

const (
	UnknownConfig ConfigType = iota
	DefaultConfig
	EnvConfig
	FileConfig
)

var configFiles = []string{
	".commitlint.yml",
	".commitlint.yaml",
	"commitlint.yml",
	"commitlint.yaml",
}

type ConfigType byte

func (c ConfigType) String() string {
	switch c {
	case DefaultConfig:
		return "Default"
	case EnvConfig:
		return "Env"
	case FileConfig:
		return "File"
	default:
		return "Unknown"
	}
}

// LookupConfigPath returns config file path following below order
//  1. env path
// 	2. commitlint.yaml in current directory
// 	3. use default config
func LookupConfigPath() (confPath string, typ ConfigType, err error) {
	envConf := os.Getenv(CommitlintConfigEnv)
	if envConf != "" {
		envConf = filepath.Clean(envConf)
		isExists, ferr := isFileExists(envConf)
		if ferr != nil {
			return "", UnknownConfig, err
		}
		if isExists {
			return envConf, EnvConfig, nil
		}
	}

	// get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", UnknownConfig, err
	}

	// check if conf file exists in current directory
	for _, confFile := range configFiles {
		currentDirConf := filepath.Join(currentDir, confFile)
		isExists, ferr := isFileExists(currentDirConf)
		if ferr != nil {
			return "", UnknownConfig, err
		}
		if isExists {
			return currentDirConf, FileConfig, nil
		}
	}

	// default config
	return "", DefaultConfig, nil
}

func isFileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return err == nil, err
}
