COMMANDS=ls

BINARIES=$(addprefix bin/,$(COMMANDS))

GOPATH=$(HOME)/kit

.PHONY: plugins export kit

FORCE:

bin/%: plugins/% FORCE
	@echo "$@"
	@go build -buildmode=plugin -o $@ ./$<

plugins: $(BINARIES)
	@echo "$@"
	@go run ./cmd/publish/main.go

kit: vendor plugins
	@echo "$@"
	@go build -o kit ./cmd/kit/main.go

vendor:
	@echo "$@"
	@go get -u github.com/whyrusleeping/gx github.com/whyrusleeping/gx-go
	@gx lock-install

define EXPORTS
export GOPATH=$(GOPATH)
export PATH=$(GOPATH)/bin:$(PATH)
endef

export EXPORTS
export:
	@echo "$$EXPORTS"
