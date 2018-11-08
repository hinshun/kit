COMMANDS=init plugin/add plugin/rm plugin/publish

BINARIES=$(addprefix bin/,$(COMMANDS))

.PHONY: clean

kit: bin
	@echo "$@"
	@go install -ldflags "$(shell go run ./cmd/linker)" ./cmd/kit

bin: $(BINARIES)

FORCE:

bin/%: core/% FORCE
	@echo "$@"
	@go build -buildmode=plugin -o $@-linux-amd64 ./$<

clean:
	@echo "$@"
	@rm -rf bin/*
