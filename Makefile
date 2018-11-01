COMMANDS=ls init plugin/add plugin/rm plugin/publish

BINARIES=$(addprefix bin/,$(COMMANDS))

.PHONY: bootstrap clean

bin: $(BINARIES)

FORCE:

bin/%: core/% FORCE
	@echo "$@"
	@go build -buildmode=plugin -o $@-linux-amd64 ./$<
	@CC=o64-clang CXX=o64-clang++ GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=plugin -o $@-darwin-amd64 ./$<
	@kit plugin publish "$@-linux-amd64,$@-darwin-amd64"

bootstrap: $(BINARIES)
	@echo "$@"
	@go run ./cmd/bootstrap/main.go

clean:
	@echo "$@"
	@rm -rf .kit bin/*
