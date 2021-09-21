PROJECT := <compnay ssh git url>/FXUBRQ-QE/Govis-CI
BINARY=govis

SRC_DIR=./
BUILD_DIR=./

GO_SOURCES=$(wildcard $(SRC_DIR)/*.go)

GO ?= go
GIT_COMMIT=`git rev-list -1 HEAD`
GOBUILD=$(GO) build
GOBUILD_FLAGS=-ldflags "-w -s -X main.GovisVersionCommitId=${GIT_COMMIT}"

GOTEST=$(GO) test

all: binaries

.PHONY: binaries
binaries: $(BUILD_DIR)/$(BINARY)

$(BUILD_DIR)/$(BINARY): $(GO_SOURCES)
	$(GOBUILD) $(GOBUILD_FLAGS) -o $(BUILD_DIR)/$@ $(SRC_DIR)

.PHONY: linux
linux:
	$(MAKE) CGO_ENABLED=0 GOOS=linux GOARCH=amd64

.PHONY: test
test:
	$(GOTEST)

.PHONY: clean
clean:
	cd $(BUILD_DIR) && rm -rf --preserve-root $(BINARY)
