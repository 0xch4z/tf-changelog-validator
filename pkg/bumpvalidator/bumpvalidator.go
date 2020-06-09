// Package bumpvalidator provides validation for semver bumps.
package bumpvalidator

import (
	"errors"
	"fmt"

	"github.com/blang/semver"
)

// VersionKind represents a kind of version.
type VersionKind string

// VersionKind constants start with Version
const (
	VersionMajor = "major"
	VersionMinor = "minor"
	VersionPatch = "patch"
)

// VersionSkippedError is an error that signals that a version was skipped.
type VersionSkippedError struct {
	Major, Minor, Patch uint64
}

// Error implements (error).Error
func (e VersionSkippedError) Error() string {
	return fmt.Sprintf("version %d.%d.%d was skipped", e.Major, e.Minor, e.Patch)
}

// RetrogressiveUpdateError is an error that signals that a version was erroneously
// decremented.
type RetrogressiveUpdateError struct {
	From, To    uint64
	VersionKind VersionKind
}

// Error implements (error).Error
func (e RetrogressiveUpdateError) Error() string {
	return fmt.Sprintf("retrogressive %s version update from %d to %d", e.VersionKind, e.From, e.To)
}

// Validate validates the updating from one version to another.
func Validate(old, new semver.Version) error {
	if new.Major == old.Major+1 {
		if new.Minor == 0 && new.Patch == 0 {
			// valid major version bump
			return nil
		} else {
			return VersionSkippedError{
				Major: new.Major,
				Minor: 0,
				Patch: 0,
			}
		}
	} else if new.Major > old.Major {
		return VersionSkippedError{
			Major: old.Major + 1,
			Minor: 0,
			Patch: 0,
		}
	} else if new.Major < old.Major {
		return RetrogressiveUpdateError{
			VersionKind: VersionMajor,
			From:        old.Major,
			To:          new.Major,
		}
	}

	if new.Minor == old.Minor+1 {
		if new.Patch == 0 {
			// valid minor version bump
			return nil
		} else {
			return VersionSkippedError{
				Major: new.Major,
				Minor: new.Minor,
				Patch: 0,
			}
		}
	} else if new.Minor > old.Minor {
		return VersionSkippedError{
			Major: new.Major,
			Minor: old.Minor + 1,
			Patch: 0,
		}
	} else if new.Minor < old.Minor {
		return RetrogressiveUpdateError{
			VersionKind: VersionMinor,
			From:        old.Minor,
			To:          new.Minor,
		}
	}

	if new.Patch == old.Patch+1 {
		// valid patch version bump
		return nil
	} else if new.Patch > old.Patch {
		return VersionSkippedError{
			Major: new.Major,
			Minor: new.Minor,
			Patch: old.Patch + 1,
		}
	} else if new.Patch < old.Patch {
		return RetrogressiveUpdateError{
			VersionKind: VersionPatch,
			From:        old.Patch,
			To:          new.Patch,
		}
	}

	return errors.New("version did not change")
}
