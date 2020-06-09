package bumpvalidator

import (
	"testing"

	"github.com/blang/semver"
	"github.com/google/go-cmp/cmp"
)

func TestValidate(t *testing.T) {
	for _, fixture := range []struct {
		name      string
		err       error
		errString string
		old, new  semver.Version
	}{
		{
			name: "successful major bump",
			old:  semver.MustParse("1.12.13"),
			new:  semver.MustParse("2.0.0"),
			err:  nil,
		},
		{
			name: "major version skipped",
			old:  semver.MustParse("1.2.3"),
			new:  semver.MustParse("3.0.0"),
			err: VersionSkippedError{
				Major: 2,
			},
		},
		{
			name: "major version retrogressive update",
			old:  semver.MustParse("1.2.3"),
			new:  semver.MustParse("0.0.0"),
			err: RetrogressiveUpdateError{
				VersionKind: VersionMajor,
				From:        1,
				To:          0,
			},
		},
		{
			name: "major retrogressive update error is prioritized over minor, patch version skipped",
			old:  semver.MustParse("2.4.5"),
			new:  semver.MustParse("1.4.5"),
			err: RetrogressiveUpdateError{
				VersionKind: VersionMajor,
				From:        2,
				To:          1,
			},
		},
		{
			name: "major version skipped prioritized over minor version skipped",
			old:  semver.MustParse("1.2.3"),
			new:  semver.MustParse("3.3.3"),
			err: VersionSkippedError{
				Major: 2,
			},
		},
		{
			name: "major version skipped prioritized over patch version skipped",
			old:  semver.MustParse("2.2.3"),
			new:  semver.MustParse("4.0.1"),
			err: VersionSkippedError{
				Major: 3,
			},
		},
		{
			name: "minor version skipped in major version bump",
			old:  semver.MustParse("3.1.1"),
			new:  semver.MustParse("4.1.0"),
			err: VersionSkippedError{
				Major: 4,
			},
		},
		{
			name: "patch version skipped in major version bump",
			old:  semver.MustParse("4.2.1"),
			new:  semver.MustParse("5.0.1"),
			err: VersionSkippedError{
				Major: 5,
			},
		},

		{
			name: "successful minor bump",
			old:  semver.MustParse("1.2.3"),
			new:  semver.MustParse("1.3.0"),
			err:  nil,
		},
		{
			name: "minor version skipped",
			old:  semver.MustParse("2.14.5"),
			new:  semver.MustParse("2.16.0"),
			err: VersionSkippedError{
				Major: 2, Minor: 15,
			},
		},
		{
			name: "minor version retrogressive update",
			old:  semver.MustParse("1.5.1"),
			new:  semver.MustParse("1.4.0"),
			err: RetrogressiveUpdateError{
				VersionKind: VersionMinor,
				From:        5,
				To:          4,
			},
		},
		{
			name: "minor retrogressive update error is prioritized over patch version skipped",
			old:  semver.MustParse("6.4.3"),
			new:  semver.MustParse("6.2.3"),
			err: RetrogressiveUpdateError{
				VersionKind: VersionMinor,
				From:        4,
				To:          2,
			},
		},
		{
			name: "minor version skipped prioritized over patch version skipped",
			old:  semver.MustParse("2.2.3"),
			new:  semver.MustParse("2.4.1"),
			err: VersionSkippedError{
				Major: 2,
				Minor: 3,
			},
		},
		{
			name: "patch version skipped in minor version patch",
			old:  semver.MustParse("2.6.1"),
			new:  semver.MustParse("2.7.1"),
			err: VersionSkippedError{
				Major: 2,
				Minor: 7,
			},
		},

		{
			name: "successful patch bump",
			old:  semver.MustParse("1.2.3"),
			new:  semver.MustParse("1.2.4"),
			err:  nil,
		},
		{
			name: "patch version skipped",
			old:  semver.MustParse("11.2.12"),
			new:  semver.MustParse("11.2.19"),
			err: VersionSkippedError{
				Major: 11, Minor: 2, Patch: 13,
			},
		},
		{
			name: "patch version retrogressize update",
			old:  semver.MustParse("0.2.3"),
			new:  semver.MustParse("0.2.2"),
			err: RetrogressiveUpdateError{
				VersionKind: VersionPatch,
				From:        3,
				To:          2,
			},
		},

		{
			name:      "no change throws error",
			old:       semver.MustParse("1.2.3"),
			new:       semver.MustParse("1.2.3"),
			errString: "version did not change",
		},
	} {
		err := Validate(fixture.old, fixture.new)
		t.Run(fixture.name, func(t *testing.T) {
			if fixture.errString != "" {
				if err == nil {
					t.Errorf("expected error to be thrown matching '%s' but got nil", fixture.errString)
				} else if err.Error() != fixture.errString {
					t.Errorf("expected error to match '%s' but got '%s'", fixture.errString, err.Error())
				}
				return
			}

			if diff := cmp.Diff(fixture.err, err); diff != "" {
				t.Errorf("error was not as expected for bump %s -> %s:\n%s", fixture.old, fixture.new, diff)
			}
		})
	}
}
