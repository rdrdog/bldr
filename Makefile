VERSION ?= $(shell git describe --tags --dirty --always | sed -e 's/^v//')
IS_SNAPSHOT = $(if $(findstring -, $(VERSION)),true,false)
MAJOR_VERSION = $(word 1, $(subst ., ,$(VERSION)))
MINOR_VERSION = $(word 2, $(subst ., ,$(VERSION)))
PATCH_VERSION = $(word 3, $(subst ., ,$(word 1,$(subst -, , $(VERSION)))))
NEW_VERSION ?= $(MAJOR_VERSION).$(MINOR_VERSION).$(shell echo $$(( $(PATCH_VERSION) + 1)) )

HAS_TOKEN = $(if $(test -e ~/.config/github/token),true,false)
ifeq (true,$(HAS_TOKEN))
	export GITHUB_TOKEN := $(shell cat ~/.config/github/token)
endif

.PHONY: pr
pr: tidy format test

.PHONY: dist
dist:
	go build -ldflags "-X main.version=$(VERSION)" -o dist/local/bldr cmd/bldr/main.go
	go build -ldflags "-X main.version=$(VERSION)" -o dist/local/dplyr cmd/dplyr/main.go

.PHONY: format
format:
	go fmt ./...

.PHONY: test
test:
	go test ./... -v -covermode=atomic

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: generate
generate:
	go generate ./...

.PHONY: promote
promote:
	@git fetch --tags
	@echo "VERSION:$(VERSION) IS_SNAPSHOT:$(IS_SNAPSHOT) NEW_VERSION:$(NEW_VERSION)"
ifeq (false,$(IS_SNAPSHOT))
	@echo "Unable to promote a non-snapshot"
	@exit 1
endif
ifneq ($(shell git status -s),)
	@echo "Unable to promote a dirty workspace"
	@exit 1
endif
	git tag -a -m "releasing v$(NEW_VERSION)" v$(NEW_VERSION)
	git push origin v$(NEW_VERSION)

.PHONY: snapshot
snapshot:
	goreleaser build \
		--rm-dist \
		--single-target \
		--snapshot
