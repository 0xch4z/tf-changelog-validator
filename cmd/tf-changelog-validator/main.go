package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Charliekenney23/tf-changelog-validator/internal/preflight"
	"github.com/Charliekenney23/tf-changelog-validator/pkg/chlogvalidator"
	"github.com/go-git/go-git/v5"
)

const (
	repoPathEnvVar = "REPO_PATH"
)

var repoPath string

func init() {
	flag.Usage = usage
	flag.StringVar(&repoPath, "repoPath", ".", "Path to the Terraform provider's git repository.")
	flag.Parse()

	if repoPathVal, ok := os.LookupEnv(repoPathEnvVar); ok {
		repoPath = repoPathVal
	}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), `Usage:    tf-changelog-validator [OPTIONS]

Validate your Terraform provider's changelog to prevent release errors

Options:
`)
	flag.PrintDefaults()
}

func main() {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("failed to open git repository '%s': %s", repoPath, err)
	}

	ok, err := preflight.Check(r)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Println("[INFO] preflight check failed; not running validation")
		return
	}

	changelogPath := filepath.Join(repoPath, "CHANGELOG.md")
	f, err := os.Open(changelogPath)
	if err != nil {
		log.Fatalf("failed to open changelog '%s': %s", changelogPath, err)
	}
	log.Printf("[INFO] validating changelog '%s'", changelogPath)

	if err := chlogvalidator.Validate(f); err != nil {
		log.Fatalf("validation failed: %s", err)
	}
	log.Printf("[INFO] changelog is valid")
}
