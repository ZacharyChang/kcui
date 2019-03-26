BINARY:=kcui
SYSTEM:=
CHECKS:=check
BUILDOPTS:=-v
LDFLAGS="-w -s"
CGO_ENABLED:=0

.PHONY: all
all: build

.PHONY: build
build: clean
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags=$(LDFLAGS) -o $(BINARY)

.PHONY: linux
linux: clean
	$(eval SYSTEM := GOOS=linux)
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags=$(LDFLAGS) -o $(BINARY)

.PHONY: test
test:
	GO111MODULE=on go test -v -race ./...

.PHONY: clean
clean:
	GO111MODULE=on go clean
	rm -f kcui

.PHONY: compress
compress: $(BINARY)
	upx $(BINARY)