IS_LINUX=$(shell sed --version > /dev/null 2> /dev/null && echo $$?)
ifeq ($(IS_LINUX),0)
	SED_IN_PLACE=-i
else
	SED_IN_PLACE=-i ""
endif


GITHASH := $(shell git rev-parse --short HEAD)
BUILDDATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

.PHONY: all
all: tidy deps build

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
build-ui: ##  build the executable and set the version
	cd ui && npm install && ng build --base-href /ui/

.PHONY: fmt
fmt:
	find . -type f -name '*.go' -not -path "./vendor/*" -exec goimports -w  {} \;


