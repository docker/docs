package config

import (
	"fmt"
	"strconv"
	"strings"
)

// Version is a major/minor version pair of the form Major.Minor
// Major version upgrades indicate structure or type changes
// Minor version upgrades should be strictly additive
type Version string

// MajorMinorVersion constructs a Version from its Major and Minor components
func MajorMinorVersion(major, minor uint) Version {
	return Version(fmt.Sprintf("%d.%d", major, minor))
}

// CurrentVersion is the most recent Version that can be parsed
var CurrentVersion = MajorMinorVersion(0, 1)

func (version Version) major() (uint, error) {
	majorPart := strings.Split(string(version), ".")[0]
	major, err := strconv.ParseUint(majorPart, 10, 0)
	return uint(major), err
}

// Major returns the major version portion of a Version
func (version Version) Major() uint {
	major, _ := version.major()
	return major
}

func (version Version) minor() (uint, error) {
	minorPart := strings.Split(string(version), ".")[1]
	minor, err := strconv.ParseUint(minorPart, 10, 0)
	return uint(minor), err
}

// Minor returns the minor version portion of a Version
func (version Version) Minor() uint {
	minor, _ := version.minor()
	return minor
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
// Unmarshals a string of the form X.Y into a Version, validating that X and Y can represent uints
func (version *Version) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var versionString string
	err := unmarshal(&versionString)
	if err != nil {
		return err
	}

	newVersion := Version(versionString)
	if _, err := newVersion.major(); err != nil {
		return err
	}

	if _, err := newVersion.minor(); err != nil {
		return err
	}

	*version = newVersion
	return nil
}
