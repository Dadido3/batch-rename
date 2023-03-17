// Copyright (c) 2021-2023 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// File for handling versioning.
// This relies on naming git tags in the semantic version scheme, and on the correct forwarding of those tag names into the build process.
// If an invalid version string is supplied, the version will fall back to some default value containing a version control hash if possible.

package main

import (
	"runtime/debug"
	"strconv"

	"github.com/Masterminds/semver/v3"
)

// version contains the semantic version of the software as a string.
//
// This variable is only used to embed the version into this program.
// To get the current version, use the `Version` variable instead.
//
// To compile the software with embedded version information, use the following command:
//
//	go build -ldflags="-X 'main.version=x.y.z'"
//
// where `x.y.z` is the desired version as SemVer string.
//
// This may or may not contain the `v` prefix.
var version = ""

// Version stores the current version of the program as a `Masterminds/semver` object.
var Version = parseVersion(version)

// getVersion returns the version of the program.
func parseVersion(version string) *semver.Version {
	// Parse version tag.
	// If it can't be parsed, fall back to the code below.
	if version != "" {
		if versionParsed, err := semver.NewVersion(version); err == nil {
			return versionParsed
		}
	}

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		var vcsRevision string
		var vcsModified bool
		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				vcsRevision = setting.Value
			case "vcs.modified":
				vcsModified, _ = strconv.ParseBool(setting.Value)
			}
		}

		// Try to use the version information from debug build info.
		// This only works when the software is installed with `go install`, as otherwise `buildInfo.Main.Version` will contain a revision hash or something else.
		// If this fails, fall back to code below.
		if versionParsed, err := semver.NewVersion(buildInfo.Main.Version); err == nil {
			return versionParsed
		}

		// Generate virtual version, as we don't know the current version.
		// Use 0.0.0, even though it's not a valid SemVer version.
		// TODO: Implement a way to retrieve the version tag, and increment the patch version
		semVer := "0.0.0"

		if vcsModified {
			// We don't know the current version tag, and there are uncommitted changes.
			semVer += "-unknown.dirty"
		} else {
			// We don't know the current version tag.
			semVer += "-unknown"
		}

		// Add VCS revision string.
		// TODO: Shorten revision string to 7 characters
		semVer += "+" + vcsRevision

		if versionParsed, err := semver.NewVersion(semVer); err == nil {
			return versionParsed
		}
	}

	return semver.MustParse("0.0.0-unknown")
}
