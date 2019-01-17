# kit

# plugin was built with a different version of package

Waiting for Go 1.13 to resolve reproducible builds issue (otherwise plugins need to have identical `go env`) in order to load properly. Without this, `kit` compiled with a different `GOPATH` cannot load plugins built by a CI server with a different `GOPATH`:
- https://github.com/golang/go/issues/26759#issuecomment-438774260
- https://github.com/golang/go/issues/16860#issuecomment-440317779
