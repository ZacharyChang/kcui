BINARY:=kcui
SYSTEM:=
CHECKS:=check
BUILDOPTS:=-v
LDFLAGS="-w -s"
CGO_ENABLED:=0
BUILD_DIR=build

.PHONY: all
all: build

.PHONY: build
build: clean
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags=$(LDFLAGS) -o $(BINARY)

.PHONY: linux
linux: clean
	$(eval SYSTEM := GOOS=linux)
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags=$(LDFLAGS) -o $(BINARY)

.PHONY: all
all: clean
	$(eval SYSTEM := GOOS=linux)
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)_linux
	$(eval SYSTEM := GOOS=darwin)
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)_darwin
	$(eval SYSTEM := GOOS=windows)
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)_windows.exe

.PHONY: test
test:
	GO111MODULE=on ginkgo -v -race --cover ./...

.PHONY: clean
clean:
	GO111MODULE=on go clean
	rm -f kcui

.PHONY: compress
compress: $(BINARY)
	upx $(BINARY)

.PHONY: compress-all
compress-all: all
	upx -9 $(BUILD_DIR)/*