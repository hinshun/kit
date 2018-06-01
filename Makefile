COMMANDS=ls

BINARIES=$(addprefix bin/,$(COMMANDS))

FORCE:

bin/%: plugins/% FORCE
	@echo "$@"
	@go build -buildmode=plugin -o $@ ./$<
	@ipfs add $@ -q | ipfs name publish

binaries: $(BINARIES)
	@echo "$@"
