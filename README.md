# kit

A composable command line toolkit distributed commands through [IPFS](github.com/ipfs/go-ipfs)

# Getting started

Before you begin, you'll need:
- IPFS daemon [go-ipfs](https://github.com/ipfs/go-ipfs#install) or [js-ipfs](https://github.com/ipfs/go-ipfs#install).
- Docker daemon [docker-ce](https://docs.docker.com/install)

We can start our IPFS daemon, and then `make` will do the following:
- Cross compile core plugins needed to get started
- Generate `ldflags` to link core plugins CID to `kit` global variables for its core plugins
- Cross-compile `kit` for `darwin-amd64` and `linux-amd64`.

```sh
$ ipfs init
initializing IPFS node at /home/edgarl/.ipfs
// ...

$ ipfs daemon
Initializing daemon...
// ...

$ make
core/init
make[1]: Entering directory '/home/edgarl/go/src/github.com/hinshun/kit'
cross
Sending build context to Docker daemon  65.61MB
```

# Current issues

## plugin doesn't use CGO but still require cross-compiling toolchain

When cross-compiling go plugins from linux to darwin or darwin to linux, it uses external linking, therefore you need to have cross compiling toolchain for the target platform installed. For example, to cross-compile go plugin from linux to darwin, you need the darwin cross compiling toolchain (binaries like `/Users/urso/.gvm/gos/go1.8beta1/pkg/tool/darwin_amd64/link`). Go could compile using internal linking if CGO isn't used but that is only a punted feature request at this moment.
- https://github.com/golang/go/issues/18157

## plugin was built with a different version of package

Waiting for Go 1.13 to resolve reproducible builds issue (otherwise plugins need to have identical `go env`) in order to load properly. Without this, `kit` compiled with a different `GOPATH` cannot load plugins built by a CI server with a different `GOPATH`:
- https://github.com/golang/go/issues/26759#issuecomment-438774260
- https://github.com/golang/go/issues/16860#issuecomment-440317779

Note: Use '-trimpath' to have more reproducible builds:
- https://github.com/github/hub/pull/1994/files
