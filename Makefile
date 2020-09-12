BINARY:=kcui
SYSTEM:=
CHECKS:=check
BUILDOPTS:=-v
PKG:=github.com/ZacharyChang/kcui
VERSION:=`git describe --tags`
GIT_HASH:=`git describe --always`
LDFLAGS="-w -s -X $(PKG)/version.Version=$(VERSION) -X $(PKG)/version.GitHash=$(GIT_HASH)"
CGO_ENABLED:=0
BUILD_DIR=build
GO111MODULE=off

.PHONY: all
all: build

.PHONY: build
build: clean
	GO111MODULE=$(GO111MODULE) CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags=$(LDFLAGS) -o $(BINARY)

.PHONY: linux
linux: clean
	$(eval SYSTEM := GOOS=linux)
	GO111MODULE=$(GO111MODULE) CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags=$(LDFLAGS) -o $(BINARY)

.PHONY: all
all: clean
	$(eval SYSTEM := GOOS=linux)
	GO111MODULE=$(GO111MODULE) CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)_linux
	$(eval SYSTEM := GOOS=darwin)
	GO111MODULE=$(GO111MODULE) CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)_darwin
	$(eval SYSTEM := GOOS=windows)
	GO111MODULE=$(GO111MODULE) CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)_windows.exe

.PHONY: test
test:
	GO111MODULE=$(GO111MODULE) ginkgo -v -race --cover ./...

.PHONY: clean
clean:
	GO111MODULE=$(GO111MODULE) go clean
	rm -f kcui

.PHONY: compress
compress: $(BINARY)
	upx $(BINARY)

.PHONY: compress-all
compress-all: all
	upx -9 $(BUILD_DIR)/$(BINARY)_linux
	upx -9 $(BUILD_DIR)/$(BINARY)_windows.exe