GO ?= go
DIST_DIR := dist
GIT_VERSION ?= $(shell ./hack/git-version.sh)

.PHONY: build
build:
	$(GO) build -o $(DIST_DIR)/kubectl-view_serviceaccount_kubeconfig cmd/kubectl-view_serviceaccount_kubeconfig.go

TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
GORELEASER_BIN := bin/goreleaser
GORELEASER := $(TOOLS_DIR)/$(GORELEASER_BIN)
GOLANGCI_LINT_BIN := bin/golangci-lint
GOLANGCI_LINT := $(TOOLS_DIR)/$(GOLANGCI_LINT_BIN)
GOMPLATE_BIN := bin/gomplate
GOMPLATE := $(TOOLS_DIR)/$(GOMPLATE_BIN)
VALIDATE_KREW_MAIFEST_BIN := bin/validate-krew-manifest
VALIDATE_KREW_MAIFEST := $(TOOLS_DIR)/$(VALIDATE_KREW_MAIFEST_BIN)

$(GORELEASER): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR) && $(GO) build -o $(GORELEASER_BIN) github.com/goreleaser/goreleaser

$(GOLANGCI_LINT): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR) && $(GO) build -o $(GOLANGCI_LINT_BIN) github.com/golangci/golangci-lint/cmd/golangci-lint

$(GOMPLATE): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR) && $(GO) build -o $(GOMPLATE_BIN) github.com/hairyhenderson/gomplate/v3/cmd/gomplate

$(VALIDATE_KREW_MAIFEST): $(TOOLS_DIR)/go.mod
	cd $(TOOLS_DIR) && $(GO) build -o $(VALIDATE_KREW_MAIFEST_BIN) sigs.k8s.io/krew/cmd/validate-krew-manifest

.PHONY: build-cross
build-cross: $(GORELEASER)
	$(GORELEASER) build --snapshot --rm-dist

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: lint
lint: vet fmt $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

.PHONY: validate-krew-manifest
validate-krew-manifest: $(VALIDATE_KREW_MAIFEST)
	$(VALIDATE_KREW_MAIFEST) -manifest dist/view-serviceaccount-kubeconfig.yaml -skip-install

.PHONY: dist
dist: $(GORELEASER) $(GOMPLATE)
	$(GORELEASER) release --rm-dist --skip-publish --snapshot
	GIT_VERSION=$(GIT_VERSION) $(GOMPLATE) -f ./hack/view-serviceaccount-kubeconfig.yaml.tmpl > $(DIST_DIR)/view-serviceaccount-kubeconfig.yaml

.PHONY: clean
clean:
	rm -rf $(DIST_DIR) $(TOOLS_BIN_DIR)
