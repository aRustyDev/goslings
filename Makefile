
.PHONY: build test

APP_NAME=goslings
VERSION = $(shell rg -No "\d+\.\d+\.\d+" internal/utils/version.go)
MAKE_DEPS= gh go git rg docker
USER=arustydev
GHCR=github.com/$(USER)/$(APP_NAME)
DOCKER_CR=docker.io/$(USER)/$(APP_NAME)

build:
	@go build -o cmd/cli/$(APP_NAME) cmd/cli/main.go
	@go build -o cmd/tui/$(APP_NAME) cmd/tui/main.go
	@go build -o cmd/api/$(APP_NAME) cmd/api/main.go
	@go build -o cmd/headless/$(APP_NAME) cmd/headless/main.go

release:
	#Building CLI
	@docker build --build-arg VERSION=$(VERSION) --build-arg NAME=$(APP_NAME) --target cli ./build/
	#Building TUI
	@docker build --build-arg VERSION=$(VERSION) --build-arg NAME=$(APP_NAME) --target tui ./build/
	#Building API
	@docker build --build-arg VERSION=$(VERSION) --build-arg NAME=$(APP_NAME) --target api ./build/
	#Building Headless
	@docker build --build-arg VERSION=$(VERSION) --build-arg NAME=$(APP_NAME) --target headless ./build/
	#Tagging
	@git tag -a $(VERSION)
	@git push origin tag $(VERSION)
	@docker image tag $(APP_NAME):$(VERSION) $(GHCR):$(VERSION)
	@docker image push $(GHCR):$(VERSION)
	@docker image tag $(APP_NAME):$(VERSION) $(DOCKER_CR):$(VERSION)
	@docker image push $(DOCKER_CR):$(VERSION)
	#Creating Release for GitHub
	@gh release create v$(VERSION) --title "v$(VERSION)" --notes-from-tag --fail-on-no-commits 'cmd/cli/$(APP_NAME)#CLI Binary' 'cmd/tui/$(APP_NAME)#TUI Binary' 'cmd/api/$(APP_NAME)#API Binary' 'cmd/headless/$(APP_NAME)#Headless Binary'
	#Building Container Images
	#Building Container Images

pre-publish:
	#Building CLI
	@docker build --build-arg VERSION=$(VERSION) --build-arg NAME=$(APP_NAME) --target cli ./build/
	#Building TUI
	@docker build --build-arg VERSION=$(VERSION) --build-arg NAME=$(APP_NAME) --target tui ./build/
	#Building API
	@docker build --build-arg VERSION=$(VERSION) --build-arg NAME=$(APP_NAME) --target api ./build/
	#Building Headless
	@docker build --build-arg VERSION=$(VERSION) --build-arg NAME=$(APP_NAME) --target headless ./build/
	#Tagging
	@git tag -a $(VERSION)
	@git push origin tag $(VERSION)
	@docker image tag $(APP_NAME):$(VERSION) $(GHCR):$(VERSION)
	@docker image push $(GHCR):$(VERSION)
	@docker image tag $(APP_NAME):$(VERSION) $(DOCKER_CR):$(VERSION)
	@docker image push $(DOCKER_CR):$(VERSION)
	#Creating Release for GitHub
	@gh release create v$(VERSION) --title "v$(VERSION)" --draft --prerelease --notes-from-tag 'cmd/cli/$(APP_NAME)#CLI Binary' 'cmd/tui/$(APP_NAME)#TUI Binary' 'cmd/api/$(APP_NAME)#API Binary' 'cmd/headless/$(APP_NAME)#Headless Binary'

test:
	@go test -v

clean:
	@rm cmd/cli/$(APP_NAME)
	@rm cmd/tui/$(APP_NAME)
	@rm cmd/api/$(APP_NAME)
	@rm cmd/headless/$(APP_NAME)
