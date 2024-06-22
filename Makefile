IS_LINUX=$(shell sed --version > /dev/null 2> /dev/null && echo $$?)
ifeq ($(IS_LINUX),0)
	SED_IN_PLACE=-i
else
	SED_IN_PLACE=-i ""
endif


GITHASH := $(shell git rev-parse --short HEAD)
BUILDDATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')




.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: all
clean: ## clean the project
	rm -rf helm-live
	rm -rf ui/dist
	rm -rf internal/server/static

.PHONY: all
all: tidy deps build

.PHONY: build
build: clean build-ui build-go  ## build the golang and the ui

.PHONY: tidy
tidy: ## get the golang dependencies in the vendor folder
	GO111MODULE=on  go mod tidy

.PHONY: deps
deps: ## get the golang dependencies in the vendor folder
	GO111MODULE=on  go mod vendor

.PHONY: build-go
build-go: ##  build the executable and set the version
	go generate ./...
	go build -o helm-live -tags=jsoniter ./cmd/live

.PHONY: build-ui
build-ui: ##  build the statics web files
	cd ui && npm install && ng build --base-href /ui/

.PHONY: fmt
fmt:
	find . -type f -name '*.go' -not -path "./vendor/*" -exec goimports -w  {} \;


