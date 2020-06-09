# tf-changelog-validator [![GoDoc][godoc-badge]][godoc] ![build][build-badge]

> Validate your Terraform provider's changelog to prevent release errors.

The main purpose of this validator is to ensure that versions are not erroneously bumped. For instance: skipping a version (`1.2.3` -> `1.2.5`) or accidentally rolling back a version (`1.2.3` -> `1.2.2`).

## Usage

Install binary with go:

```
go install github.com/Charliekenney23/tf-changelog-validator/cmd/tf-changelog-validator
```

CLI usage:

```
❯ tf-changelog-validator -h                    
Usage:    tf-changelog-validator [OPTIONS]

Validate your Terraform provider's changelog to prevent release errors

Options:
  -repoPath string
    	Path to the Terraform provider's git repository. (default ".")
```

Run in a docker container:

```
docker run -v `pwd`:/var/repo:ro charliekenney23/tf-changelog-validator
```

Programmatic usage:

```golang
import "os"
import "github.com/Charliekenney23/tf-changelog-validator/pkg/chlogvalidator"

func main() {
    changelogFile, err := os.Open("./Changelog.md")
    if err != nil {
        panic(err)
    }

    if err := chlogvalidator.Validate(changelogFile); err != nil {
        log.Printf("validation error: %s", err)
    }
}
```

[build-badge]: https://github.com/Charliekenney23/tf-changelog-validator/workflows/build/badge.svg
[godoc-badge]: https://godoc.org/github.com/Charliekenney23/tf-changelog-validator?status.svg
[godoc]: https://godoc.org/github.com/Charliekenney23/tf-changelog-validator
