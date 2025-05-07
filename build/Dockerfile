# -----|Builds from Src|------------
# syntax=docker/dockerfile:1
FROM golang:1.23 AS build-api
WORKDIR /src
ARG VERSION
COPY . /src/
RUN go build -ldflags "-X main.version=$VERSION" -o cmd/goslings cmd/api/main.go

FROM golang:1.23 AS build-cli
WORKDIR /src
ARG VERSION
ARG NAME
COPY . /src/
RUN go build -ldflags "-X main.version=$VERSION" -o cmd/goslings cmd/cli/main.go

FROM golang:1.23 AS build-headless
WORKDIR /src
ARG VERSION
ARG NAME
COPY . /src/
RUN go build -ldflags "-X main.version=$VERSION" -o cmd/goslings cmd/headless/main.go

FROM golang:1.23 AS build-tui
WORKDIR /src
ARG VERSION
ARG NAME
COPY . /src/
RUN go build -ldflags "-X main.version=$VERSION" -o cmd/goslings cmd/tui/main.go

# -----|Binary Distros|------------

FROM scratch AS api
ARG VERSION
LABEL version=v$VERSION
LABEL ui="api"
LABEL description="The Untitled Goose Tool rewritten into Go, \
    this is the cloud-native version. This is \
    intended to run in a kubernetes cluster \
    and has an associated helm chart."
COPY --from=build-api cmd/goslings /bin/goslings
CMD ["/bin/goslings"]

FROM scratch AS headless
ARG VERSION
LABEL version=v$VERSION
LABEL ui="headless"
LABEL description="The Untitled Goose Tool rewritten into Go, \
    This version is a 'fire & forget' version, it \
    does not take input after it is started, it \
    can only be killed."
COPY --from=build-headless cmd/goslings /bin/goslings
CMD ["/bin/goslings"]

FROM scratch AS cli
ARG VERSION
LABEL version=v$VERSION
LABEL ui="cli"
LABEL description="The Untitled Goose Tool rewritten into Go, \
    this version is intended to be run from a \
    terminal or some CI/CD fashion. It can be used \
    both interactively and in scripted fashion."
COPY --from=build-cli cmd/goslings /bin/goslings
CMD ["/bin/goslings"]

FROM scratch AS tui
ARG VERSION
LABEL version=v$VERSION
LABEL ui="tui"
LABEL description="The Untitled Goose Tool rewritten into Go, \
    This version is strictly for interactive use, \
    it is built to allow users to have tight \
    control over the tool, and is useful for \
    debugging the tool or exploring data."
COPY --from=build-tui cmd/goslings /bin/goslings
CMD ["/bin/goslings"]
