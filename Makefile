COMMANDS=ls init plugin/add plugin/rm

BINARIES=$(addprefix bin/,$(COMMANDS))

.PHONY: bootstrap clean

bin: $(BINARIES)

FORCE:

bin/%: core/% FORCE
	@echo "$@"
	@go build -buildmode=plugin -o $@ ./$<

bootstrap: $(BINARIES)
	@echo "$@"
	@go run ./cmd/bootstrap/main.go

clean:
	@echo "$@"
	@rm -rf .kit bin/*
