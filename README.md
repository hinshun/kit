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

You will then have access to `kit init`, which is a bootstrapping case, so that if you ever end up at an empty state, you can run it to reinitialize yourself.

```
$ kit
Usage:
  kit - Composable command-line toolkit.
  
  kit [global options] command [options] <arguments>
    
Commands:
  init
      Initializes a kit config.

$ kit init

$ kit
Usage:
  kit - Composable command-line toolkit.
  
  kit [global options] command [options] <arguments>
    
Commands:
  plugin
      Manage kit plugins.

# Add the `plugin/add` plugin as `foobar`
$ kit plugin add /foobar QmVrAu9VUahsA89Xfyuc5A3fEHWtogfdQX4MefdSicdKtJ

$ kit
Usage:
  kit - Composable command-line toolkit.

  kit [global options] command [options] <arguments>

Commands:
  foobar [--usage <string>] [--pin] [--overwrite] [--] <command path> <manifest>
    Adds a plugin to kit.
    [--usage <string>]: Specify usage help text for the plugin.
    [--pin]: Pins the plugin's parent namespace if adding to an implicit namespace.
    [--overwrite]: Overwrites any namespace or command if conflicting at the command path.
    <command path>: The command path to add the plugin.
    <manifest>: The content address or resolvable name for a plugin's metadata.

  plugin
    Manage kit plugins.

# Add an empty plugin namespace for `/some/namespace`
$ kit plugin add /some/namespace ""
    
$ kit some namespace
Usage:
  kit some namespace - A plugin namespace.

  kit [global options] command [options] <arguments>

Commands:
  No commands in kit some namespace.    
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
