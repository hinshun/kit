COMMANDS=ls init plugin/add plugin/rm plugin/publish

BINARIES=$(addprefix bin/,$(COMMANDS))

.PHONY: bootstrap clean

bin: $(BINARIES)

FORCE:

bin/%: core/% FORCE
	@echo "$@"
	@go build -buildmode=plugin -o $@-linux-amd64 ./$<
	@go run ./cmd/publish "$@-linux-amd64"

bootstrap: $(BINARIES)
	@echo "$@"
	@go run ./cmd/bootstrap/main.go

clean:
	@echo "$@"
	@rm -rf .kit bin/*
