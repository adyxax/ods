SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
.DEFAULT_GOAL := help

##### Quality ##################################################################
.PHONY: check
check: ## run all code checks
	go mod verify
	go vet ./...

.PHONY: tidy
tidy: ## tidy up the code
	go fmt ./...
	go mod tidy -v

##### Development ##############################################################
.PHONY: build
build: git-crypt-unlocked ## build the code
	CGO_ENABLED=0 go build -o ./ods ./

.PHONY: clean
clean: ## clean the code
	rm -f ./ods

.PHONY: run
run: git-crypt-unlocked ## run the code
	go run ./

##### Operations ###############################################################
.PHONY: push
push: tidy no-dirty check ## push changes to git remote
	git push git main

.PHONY: deploy
deploy: build ## deploy changes to the production server
	rsync ./ods root@ods.adyxax.org:/usr/local/bin/
	ssh root@ods.adyxax.org "systemctl restart ods"

##### Utils ####################################################################
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: git-crypt-unlocked
git-crypt-unlocked:
	@git config --local --get filter.git-crypt.smudge >/dev/null || FAIL=1
	if [[ "$${FAIL:-0}" -gt 0 ]]; then
	    echo "Please unlock git-crypt before continuing!"
	    exit 1
	fi

.PHONY: help
help:
	@grep -E '^[a-zA-Z\/_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' \
	| sort

.PHONY: no-dirty
no-dirty:
	git diff --exit-code
