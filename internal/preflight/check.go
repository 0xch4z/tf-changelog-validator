package preflight

import (
	"fmt"
	"log"

	"github.com/go-git/go-git/v5"
)

const releaseBotAuthor = "tf-release-bot"

// conditionFunc returns a boolean describing whether the the condition
// passed, and optionally, an error if one occurs.
type conditionFunc func(*git.Repository) (bool, error)

// Check runs preflight checks on the Git repository to see if the
// validator should be run or not.
func Check(repo *git.Repository) (bool, error) {
	for name, conditionFunc := range map[string]conditionFunc{
		"hasTags":                           checkHasTags,
		"latestCommitAuthorIsNotReleaseBot": checkLatestCommitAuthorIsNotReleaseBot,
	} {
		ok, err := conditionFunc(repo)
		if err != nil {
			return false, fmt.Errorf("failed to evaluate condition '%s': %w", name, err)
		}
		if !ok {
			log.Printf("[INFO] preflight condition '%s' failed\n", name)
			return false, nil
		}
	}
	return true, nil
}

func checkHasTags(repo *git.Repository) (bool, error) {
	iter, err := repo.TagObjects()
	if err != nil {
		return false, fmt.Errorf("failed to get tags: %w", err)
	}
	defer iter.Close()

	tag, err := iter.Next()
	if err != nil && err.Error() != "EOF" {
		return false, fmt.Errorf("failed to get latest tag: %w", err)
	}

	if tag == nil {
		return false, nil
	}
	return true, nil
}

func checkLatestCommitAuthorIsNotReleaseBot(repo *git.Repository) (bool, error) {
	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get commits: %w", err)
	}
	defer iter.Close()

	commit, err := iter.Next()
	if err != nil {
		return false, fmt.Errorf("failed to get latest commit: %w", err)
	}

	return commit.Author.Name != releaseBotAuthor, nil
}
