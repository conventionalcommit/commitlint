package config

import (
	"fmt"
	"runtime/debug"
)

// Build constants
// all variables are set during build
var (
	version   string
	commit    string
	buildTime string
)

// Version returns short version number of the commitlint
func Version() string {
	return formShortVersion()
}

// FullVersion returns version number with hash and build time of the commitlint
func FullVersion() string {
	return formFullVersion()
}

func formShortVersion() string {
	if buildTime != "" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "v0"
	}
	return info.Main.Version
}

func formFullVersion() string {
	versionTmpl := "%s - built from %s on %s"

	if buildTime != "" {
		return fmt.Sprintf(versionTmpl, version, commit, buildTime)
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Sprintf(versionTmpl, "master", "unknown", "unknown")
	}

	var commitInfo string
	if info.Main.Sum == "" {
		commitInfo = "(" + "checksum: unknown)"
	} else {
		commitInfo = "(" + "checksum: " + info.Main.Sum + ")"
	}
	return fmt.Sprintf(versionTmpl, info.Main.Version, commitInfo, "unknown")
}
