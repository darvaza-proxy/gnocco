.PHONY: tidy tidy-root
tidy: tidy-root

tidy-root: | fmt $(REVIVE) ; $(info $(M) tidy: root)
	$(Q) $(GO) mod tidy
	$(Q) $(GO) vet ./...
	$(Q) $(REVIVE) $(REVIVE_RUN_ARGS) ./...

.PHONY: get get-root
get: get-root

get-root: ; $(info $(M) get: root)
	$(Q) $(GO) get -tags tools -v ./...
	$(Q) for url in $(GO_INSTALL_URLS); do $(GO) install -v $$url; done

.PHONY: build build-root
build: build-root

build-root: ; $(info $(M) build: root)
	$(Q) if $(GO) list ./... | grep -e '.*/cmd/[^/]\+$$' > /dev/null; then $(GO_BUILD_CMD) ./...; else $(GO_BUILD) ./...; fi

.PHONY: test test-root
test: test-root

test-root: ; $(info $(M) test: root)
	$(Q) $(GO) test ./...

.PHONY: up up-root
up: up-root

up-root: ; $(info $(M) up: root)
	$(Q) $(GO) get -tags tools -u -v ./...
	$(Q) $(GO) mod tidy
	$(Q) for url in $(GO_INSTALL_URLS); do $(GO) install -v $$url; done

root: get-root build-root tidy-root
