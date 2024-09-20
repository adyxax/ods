SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
.DEFAULT_GOAL := help
OUTDIR := "."

CONTAINER_REGISTRY ?= localhost
CONTAINER_TAG ?= latest

##### Utils ####################################################################
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: help
help:
	@grep -E '^[a-zA-Z\/_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' \
	| sort

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

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
build: ## build the code
	CGO_ENABLED=0 go build -o $(OUTDIR)/ods ./

.PHONY: clean
clean: ## clean the code
	rm -f $(OUTDIR)/ods

.PHONY: run
run: ## run the code
	go run ./

##### Containers ###############################################################
.PHONY: container-build
container-build: ## build the container image
	printf "Dockerfile GNUmakefile go.mod go.sum index.html main.go ods.txt" | xargs shasum >checksums
	podman build \
	    -v $$PWD:/usr/src/app \
	    -t $(CONTAINER_REGISTRY)/ods:$(CONTAINER_TAG) \
	    .

.PHONY: container-push
container-push: ## push the container image to the container registry
	podman push \
            $(CONTAINER_REGISTRY)/ods:$(CONTAINER_TAG)

.PHONY: container-run
container-run: ## run the code inside podman
	podman run --rm -ti \
	    -p 8090:8090 \
	    $(CONTAINER_REGISTRY)/ods:$(CONTAINER_TAG)

##### Operations ###############################################################
.PHONY: push
push: tidy no-dirty check ## push changes to git remote
	git push git master

.PHONY: deploy
deploy: ## deploy changes to the production server
	rsync $(OUTDIR)/ods root@ods.adyxax.org:/srv/ods/
	ssh root@ods.adyxax.org "systemctl restart ods"
