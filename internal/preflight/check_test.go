package preflight

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
)

type checkFixture struct {
	name   string
	passed bool
	repo   *git.Repository
}

func runTestFixture(t *testing.T, fixture checkFixture) {
	t.Helper()

	ok, err := Check(fixture.repo)
	if err != nil {
		t.Error(err)
	}

	if fixture.passed && !ok {
		t.Error("expected checks to fail but they passed.")
	} else if !fixture.passed && ok {
		t.Error("expected checks to pass but they failed.")
	}
}

func TestCheck(t *testing.T) {
	repo, err := makeTestRepoWithCommit("not-a-bot", "test", true)
	if err != nil {
		t.Error(err)
	}
	runTestFixture(t, checkFixture{
		name:   "passes with non-bot HEAD commit with tags",
		repo:   repo,
		passed: true,
	})

	repo, err = makeTestRepoWithCommit("not-a-bot", "test", false)
	if err != nil {
		t.Error(err)
	}
	runTestFixture(t, checkFixture{
		name:   "fails with non-bot HEAD commit and no tags",
		repo:   repo,
		passed: false,
	})

	repo, err = makeTestRepoWithCommit(releaseBotAuthor, "test", true)
	if err != nil {
		t.Error(err)
	}
	runTestFixture(t, checkFixture{
		name:   "failes with bot HEAD commit and tags",
		repo:   repo,
		passed: false,
	})
}

func makeTestRepoWithCommit(author, msg string, tag bool) (*git.Repository, error) {
	fs := memfs.New()
	r, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	if err := util.WriteFile(fs, "foo", []byte("foo"), 0644); err != nil {
		return nil, err
	}

	if _, err := w.Add("foo"); err != nil {
		return nil, err
	}

	sig := &object.Signature{
		Name:  author,
		Email: "test@test.com",
		When:  time.Now(),
	}

	commitHash, err := w.Commit(msg+"\n", &git.CommitOptions{
		Author: sig,
	})
	if err != nil {
		return nil, err
	}

	if tag {
		if _, err := r.CreateTag("0.0.1", commitHash, &git.CreateTagOptions{
			Tagger:  sig,
			Message: "test",
		}); err != nil {
			return nil, err
		}
	}

	return r, nil
}
