COMMANDS=ls

BINARIES=$(addprefix bin/,$(COMMANDS))

GOPATH=$(HOME)/kit

FORCE:

bin/%: plugins/% FORCE
	@echo "$@"
	@go build -buildmode=plugin -o $@ ./$<
	# @ipfs add $@ -q | ipfs name publish

binaries: $(BINARIES)
	@echo "$@"

gx:
	@echo "$@"
	@go get -u github.com/whyrusleeping/gx github.com/whyrusleeping/gx-go

define EXPORTS
export GOPATH=$(GOPATH)
export PATH=$(GOPATH)/bin:$(PATH)
endef

export EXPORTS
export:
	@echo "$$EXPORTS"
