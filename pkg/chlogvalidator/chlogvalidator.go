package chlogvalidator

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"

	"github.com/Charliekenney23/tf-changelog-validator/pkg/bumpvalidator"
	"github.com/blang/semver"
)

const (
	tagPattern             = "\\d+\\.\\d+\\.\\d+"
	entryPattern           = "## " + tagPattern
	releasedEntryPattern   = entryPattern + " \\([A-Z][a-z]+ \\d{1,2}, \\d{4}\\)"
	unreleasedEntryPattern = entryPattern + " \\(Unreleased\\)"
)

var (
	tagRe             = regexp.MustCompile(tagPattern)
	releasedEntryRe   = regexp.MustCompile(releasedEntryPattern)
	unreleasedEntryRe = regexp.MustCompile(unreleasedEntryPattern)
)

func parseLastReleasedSemver(changelogBytes []byte) (*semver.Version, error) {
	entryBytes := releasedEntryRe.Find(changelogBytes)
	if entryBytes == nil {
		return nil, nil
	}
	return parseSemverFromChangelogEntry(entryBytes)
}

func parseUnreleasedSemver(changelogBytes []byte) (*semver.Version, error) {
	entryBytes := unreleasedEntryRe.Find(changelogBytes)
	if entryBytes == nil {
		return nil, nil
	}
	return parseSemverFromChangelogEntry(entryBytes)
}

func parseSemverFromChangelogEntry(entryBytes []byte) (*semver.Version, error) {
	tagBytes := tagRe.Find(entryBytes)
	if tagBytes == nil {
		return nil, fmt.Errorf("could not find tag in entry: '%s'", string(entryBytes))
	}

	tag := string(tagBytes)
	version, err := semver.Parse(tag)
	if err != nil {
		err = fmt.Errorf("failed to parse semver from tag '%s': %s", tag, err)
	}
	return &version, err
}

// Validate validates a changelog.
func Validate(r io.Reader) error {
	changelogBytes, readErr := ioutil.ReadAll(r)
	if readErr != nil {
		return readErr
	}

	releasedVer, err := parseLastReleasedSemver(changelogBytes)
	if err != nil {
		return err
	}
	if releasedVer == nil {
		return errors.New("could not find previously released changelog entry")
	}

	unreleasedVer, err := parseUnreleasedSemver(changelogBytes)
	if err != nil {
		return err
	}
	if unreleasedVer == nil {
		return errors.New("could not find unreleased changelog entry")
	}

	return bumpvalidator.Validate(*releasedVer, *unreleasedVer)
}
