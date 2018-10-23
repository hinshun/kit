COMMANDS=ls init config/add config/rm

BINARIES=$(addprefix bin/,$(COMMANDS))

.PHONY: kit bootstrap clean

kit: bootstrap
	@echo "$@"
	@go build -o kit ./cmd/kit/main.go

FORCE:

bin/%: core/% FORCE
	@echo "$@"
	@go build -buildmode=plugin -o $@ ./$<

bin: $(BINARIES)

bootstrap: $(BINARIES)
	@echo "$@"
	@go run ./cmd/bootstrap/main.go

clean:
	@echo "$@"
	@rm -rf .kit bin/*
