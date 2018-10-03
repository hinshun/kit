COMMANDS=ls

BINARIES=$(addprefix bin/,$(COMMANDS))

GOPATH=$(HOME)/kit

.PHONY: plugins export kit clean

kit: vendor plugins
	@echo "$@"
	@go build -o kit ./cmd/kit/main.go

FORCE:

bin/%: plugins/% FORCE
	@echo "$@"
	@go build -buildmode=plugin -o $@ ./$<

vendor:
	@echo "$@"
	@go get -u github.com/whyrusleeping/gx github.com/whyrusleeping/gx-go
	@gx lock-install

plugins: $(BINARIES)
	@echo "$@"
	@go run ./cmd/publish/main.go

clean:
	@echo "$@"
	@rm -rf .kit bin/*

define EXPORTS
export GOPATH=$(GOPATH)
export PATH=$(GOPATH)/bin:$(PATH)
endef

export EXPORTS
export:
	@echo "$$EXPORTS"
