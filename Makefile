REG ?= quay.io
ORG ?= kiagnose
CHECKUP_IMAGE_NAME ?= kubevirt-dpdk-checkup
CHECKUP_IMAGE_TAG ?= latest
GO_IMAGE_NAME := docker.io/library/golang
GO_IMAGE_TAG := 1.19
BIN_DIR = $(CURDIR)/_output/bin
CRI_BIN ?= $(shell hack/detect_cri.sh)
CRI_BUILD_BASE_IMAGE_TAG ?= latest

build:
	$(CRI_BIN) run --rm \
	           --volume `pwd`:$(CURDIR):Z \
	           --workdir $(CURDIR) \
	           --user $(shell id -u):$(shell id -g) \
	           -e XDG_CACHE_HOME=/tmp/.cache \
	           -e GOOS=linux \
	           -e GOARCH=amd64 \
	           $(GO_IMAGE_NAME):$(GO_IMAGE_TAG) go build -v -o $(BIN_DIR)/$(CHECKUP_IMAGE_NAME) ./cmd/
	$(CRI_BIN) build --build-arg BASE_IMAGE_TAG=$(CRI_BUILD_BASE_IMAGE_TAG) . -t $(REG)/$(ORG)/$(CHECKUP_IMAGE_NAME):$(CHECKUP_IMAGE_TAG)
.PHONY: build
