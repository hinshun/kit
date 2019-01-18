COMMANDS=init plugin/add plugin/rm plugin/publish

CORE=$(addprefix core/,$(COMMANDS))

.PHONY: all local mod kit plugins cross clean

all: plugins kit

local:
	@echo "@"
	@make kit GATEWAY=127.0.0.1

mod:
	@echo "$@"
	@docker build -t mod -f dockerfiles/kit/Dockerfile --target mod .

kit: mod
	@echo "$@"
	@make cross PKG="./cmd/kit" BUILDMODE="default" LDFLAGS="-X github.com/hinshun/kit/content/ipfsstore.Gateway=$(GATEWAY) $$(docker run --rm -it -v $$(pwd):/src -w /src --network host mod go run ./cmd/linker /ip4/$(GATEWAY)/tcp/5001)"

plugins: $(CORE)

FORCE:

core/%: FORCE
	@echo "$@"
	@make cross PKG="./$@" BUILDMODE="plugin"

cross:
	@echo "$@"
	@docker build -t kit -f dockerfiles/kit/Dockerfile --target build --build-arg PKG="$(PKG)" --build-arg BUILDMODE="$(BUILDMODE)" --build-arg LDFLAGS="$(LDFLAGS)" .
	@docker rm kit || true
	@docker create --name kit kit bash
	@docker cp kit:/root/go/bin/. bin
	@docker rm kit

clean:
	@echo "$@"
	@rm -rf bin/*
