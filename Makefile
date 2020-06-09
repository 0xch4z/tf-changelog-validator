PROJECT := charliekenney23/tf-changelog-validator

CMD_DIR     := cmd
SCRIPTS_DIR := scripts

MAIN_PKG := $(CMD_DIR)/tf-changelog-validator

fmtcheck:
	bash $(SCRIPTS_DIR)/verify-gofmt.sh

tidycheck:
	bash $(SCRIPTS_DIR)/verify-gomod-tidy.sh

verify: fmtcheck tidycheck

test:
	go test -v ./...

install:
	(cd $(MAIN_PKG); go install)

docker-build:
	docker build . -t $(PROJECT)
