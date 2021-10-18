package config

import (
	"fmt"
	"runtime/debug"

	"golang.org/x/mod/semver"
)

// Build constants
// all variables are set during build
var (
	version   string
	commit    string
	buildTime string
)

const (
	unknownVersion   = "v0"
	unknownBuildTime = "unknown"
	unknownHash      = "unknown"
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
		return unknownVersion
	}

	if semver.IsValid(info.Main.Version) {
		return info.Main.Version
	}
	return unknownVersion
}

func formFullVersion() string {
	versionTmpl := "%s - built from %s on %s"

	if buildTime != "" {
		return fmt.Sprintf(versionTmpl, version, commit, buildTime)
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Sprintf(versionTmpl, unknownVersion, unknownHash, unknownBuildTime)
	}

	var commitInfo string
	if info.Main.Sum == "" {
		commitInfo = "(" + "checksum: " + unknownHash + ")"
	} else {
		commitInfo = "(" + "checksum: " + info.Main.Sum + ")"
	}

	var versionInfo string
	if semver.IsValid(info.Main.Version) {
		versionInfo = info.Main.Version
	} else {
		versionInfo = unknownVersion
	}
	return fmt.Sprintf(versionTmpl, versionInfo, commitInfo, unknownBuildTime)
}
