package chlogvalidator

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	for _, fixture := range []struct {
		name      string
		changelog string
		errString string
	}{
		{
			name: "passes valid major version bump",
			changelog: `
## 3.0.0 (Unreleased)

lorem ipsum.

## 2.4.1 (April 20, 2019)

lorem ipsum.

## 2.4.0 (April 19, 2019)

lorem ipsum`,
		},
		{
			name: "passes valid minor version bump",
			changelog: `
## 1.12.0 (Unreleased)
## 1.11.2 (June 05, 2020)

lorem ipsum.`,
		},
		{
			name: "passes valid patch version bump",
			changelog: `
## 4.2.12 (Unreleased)

Bug Fixes: fixed a bug

## 4.2.11 (June 8, 2015)`,
		},

		{
			name:      "fails on retrogressive update",
			errString: "retrogressive minor version update from 61 to 60",
			changelog: `
## 4.60.12 (Unreleased)
## 4.61.11 (June 24, 2020)`,
		},
		{
			name:      "fails on version skipped",
			errString: "version 5.9.0 was skipped",
			changelog: `
## 5.10.0 (Unreleased)
## 5.8.0 (November 23, 2019)`,
		},
		{
			name:      "fails when cannot find previously released entry",
			errString: "could not find previously released changelog entry",
			changelog: `
## 0.1.0 (Unreleased)
lorem ipsum.`,
		},
		{
			name:      "fails when cannot find unreleased entry",
			errString: "could not find unreleased changelog entry",
			changelog: `
## 11.2.0 (June 23, 2020)
lorem ipsum.
## 11.1.4 (June 22, 2020)
lorem ipsum.`,
		},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			err := Validate(strings.NewReader(fixture.changelog))
			if fixture.errString != "" {
				if err == nil {
					t.Errorf("expected error to be thrown matching:\n\t'%s'\nbut got:\n\tnil", fixture.errString)
				} else if err.Error() != fixture.errString {
					t.Errorf("expected error to match:\n\t'%s'\nbut got:\n\t'%s'", fixture.errString, err.Error())
				}
				return
			} else if fixture.errString == "" && err != nil {
				t.Errorf("unexpected error validating changelog:\n\t%s", err)
			}
		})
	}
}
