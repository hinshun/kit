FROM cross AS mod
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM cross AS build
ARG PKG=.
ARG BUILDMODE=default
ARG LDFLAGS=
COPY --from=mod $GOCACHE $GOCACHE
COPY --from=mod $GOPATH/pkg/mod $GOPATH/pkg/mod
WORKDIR /src
COPY . .
RUN /cross/build $PKG
