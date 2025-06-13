.PHONY: build test clean all

# ==============================
# Variables
# ==============================
APP_NAME = goslings
VERSION = $(shell rg -No "\d+\.\d+\.\d+" internal/about/version.go)
MAKE_DEPS = gh go git rg docker codecov gitsign cosign codecov-cli
USER = arustydev
GHCR = github.com/$(USER)/$(APP_NAME)
DOCKER_CR = docker.io/$(USER)/$(APP_NAME)
UI_TYPES = cli tui api headless

export APP_NAME
export VERSION
export MAKE_DEPS
export USER
export GHCR
export DOCKER_CR
export UI_TYPES

# ==============================
# Functions
# ==============================
docker_push = docker image tag $(APP_NAME):$(VERSION) $(1):$(VERSION) && docker image push $(1):$(VERSION)
update_issue = $(if $(ISSUE),gh issue comment -e $(ISSUE),echo "ERROR: need to specify ISSUE, example: make update-issue ISSUE='X'")

# ==============================
# Targets
# ==============================
build:
	@RELEASE=false ./.scripts/build.sh $(UI_TYPES)

release: tag
	#Creating Release for GitHub
	@gh release create v$(VERSION) --title "v$(VERSION)" --notes-from-tag --fail-on-no-commits 'cmd/cli/$(APP_NAME)#CLI Binary' 'cmd/tui/$(APP_NAME)#TUI Binary' 'cmd/api/$(APP_NAME)#API Binary' 'cmd/headless/$(APP_NAME)#Headless Binary'

pre-publish: tag
	#Creating Release for GitHub
	@gh release create v$(VERSION) --title "v$(VERSION)" --draft --prerelease --notes-from-tag 'cmd/cli/$(APP_NAME)#CLI Binary' 'cmd/tui/$(APP_NAME)#TUI Binary' 'cmd/api/$(APP_NAME)#API Binary' 'cmd/headless/$(APP_NAME)#Headless Binary'

build-release:
	@RELEASE=true ./.scripts/build.sh $(UI_TYPES)

tag: build-release
	@git tag -a $(VERSION)
	@git push origin tag $(VERSION)
	@$(call docker_push, $(DOCKER_CR))
	@$(call docker_push, $(GHCR))

issue:
	@gh issue create -a "@me" -e

update-issue:
	@$(call update_issue)

# https://awesome-go.com/testing/
test:
	@#mockgen -source=$(GOPATH)/pkg/mod/github.com/Azure/azure-sdk-for-go/sdk/azcore@vX.Y.Z/policy/policy.go -destination=mocks/mock_azcore_policy.go -package=mocks TokenCredential
	@#go test -race -covermode=atomic -coverprofile=test/coverage/$(VERSION).cov -v ./...
	@#go test -c -cover -o app.test && ./app.test -test.run="IntegrationTest" -test.coverprofile=integration.cov && go tool cover -html=integration.cov
	@RELEASE="false" ./.scripts/test.sh

# https://github.com/ultraware/golang-fuzz
# https://go.dev/doc/tutorial/fuzz
# https://github.com/CodeIntelligenceTesting/gofuzz
# https://go.dev/doc/security/fuzz/
# https://github.com/kisielk/goflamegraph
# https://hackernoon.com/go-the-complete-guide-to-profiling-your-code-h51r3waz
# https://adihumara.gitbooks.io/golang/content/testing/profiling.html
fuzz:
	@go test -fuzz ./...
	@go test -benchmem -cpuprofile cpu.prof -memprofile mem.prof -blockprofile block.prof -bench .
	@docker run uber/go-torch -u http://<host ip>:8080/debug/pprof -p -t=30 > torch.svg
	@echo "GET http://localhost:8080/ping" | vegeta attack -rate 250 -duration=60s | vegeta report
	@open -a `Google Chrome` torch.svg

upgrade-make-deps:
	@echo "TODO: Script upgrading make deps"
	@#brew update
	@#brew upgrade

install-make-deps:
	@echo "TODO: Script installing make deps"
	@#pip install codecov-cli --break-system-packages
	@#brew install ripgrep gh git go orbstack golangci-lint goimports

clean:
	@./scripts/clean.sh $(UI_TYPES)
