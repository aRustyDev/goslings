
.PHONY=build test

APP_NAME=goslings
VERSION = $(shell rg -No "\d+\.\d+\.\d+" internal/utils/version.go)

build:
	@go build -o cmd/cli/$(APP_NAME) cmd/cli/main.go
	@go build -o cmd/tui/$(APP_NAME) cmd/tui/main.go
	@go build -o cmd/api/$(APP_NAME) cmd/api/main.go
	@go build -o cmd/headless/$(APP_NAME) cmd/headless/main.go

release:
	@go build -ldflags "-X main.version=$(VERSION)" -o cmd/cli/$(APP_NAME) cmd/cli/main.go
	@go build -ldflags "-X main.version=$(VERSION)" -o cmd/tui/$(APP_NAME) cmd/tui/main.go
	@go build -ldflags "-X main.version=$(VERSION)" -o cmd/api/$(APP_NAME) cmd/api/main.go
	@go build -ldflags "-X main.version=$(VERSION)" -o cmd/headless/$(APP_NAME) cmd/headless/main.go
	@git tag -a $(VERSION)
	@git push origin tag $(VERSION)

test:
	@go test -v

clean:
	@rm cmd/cli/$(APP_NAME)
	@rm cmd/tui/$(APP_NAME)
	@rm cmd/api/$(APP_NAME)
	@rm cmd/headless/$(APP_NAME)
