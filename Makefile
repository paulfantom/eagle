all: check-license build generate test

GO111MODULE=on
export GO111MODULE

GITHUB_URL=github.com/paulfantom/eagle
GOOS?=$(shell uname -s | tr A-Z a-z)
GOARCH?=$(shell go env GOARCH)
OUT_DIR=_output
BIN?=eagle
VERSION?=$(shell git rev-parse --short HEAD)
DOCKER_REPO?=quay.io/paulfantom/eagle

ALL_ARCH=amd64 arm arm64 ppc64le s390x
ALL_PLATFORMS=$(addprefix linux/,$(ALL_ARCH))
ALL_BINARIES ?= $(addprefix $(OUT_DIR)/$(BIN)-, \
				$(addprefix linux-,$(ALL_ARCH)))

all: clean build

clean:
	rm -rf $(OUT_DIR)
	rm -f manifests-tool

crossbuild: $(ALL_BINARIES)

$(OUT_DIR)/$(BIN): $(OUT_DIR)/$(BIN)-$(GOOS)-$(GOARCH)
	cp $(OUT_DIR)/$(BIN)-$(GOOS)-$(GOARCH) $(OUT_DIR)/$(BIN)

$(OUT_DIR)/$(BIN)-%:
	@echo ">> building for $(GOOS)/$(GOARCH) to $(OUT_DIR)/$(BIN)-$*"
	GOARCH=$(word 2,$(subst -, ,$(*:.exe=))) \
	GOOS=$(word 1,$(subst -, ,$(*:.exe=))) \
	CGO_ENABLED=0 \
	go build --installsuffix cgo -o $(OUT_DIR)/$(BIN)-$* $(GITHUB_URL)

build: $(OUT_DIR)/$(BIN)

container: $(OUT_DIR)/$(BIN)-$(GOOS)-$(GOARCH) Dockerfile
	docker build --build-arg BINARY=$(BIN)-$(GOOS)-$(GOARCH) -t $(DOCKER_REPO):$(VERSION)-$(GOARCH) .
ifeq ($(GOARCH), amd64)
	docker tag $(DOCKER_REPO):$(VERSION)-$(GOARCH) $(DOCKER_REPO):$(VERSION)
endif


manifest-tool:
	curl -fsSL https://github.com/estesp/manifest-tool/releases/download/v1.0.2/manifest-tool-linux-amd64 > ./manifest-tool
	chmod +x ./manifest-tool

push-%:
	$(MAKE) GOARCH=$* container
	docker push $(DOCKER_REPO):$(VERSION)-$*

comma:= ,
empty:=
space:= $(empty) $(empty)
manifest-push: manifest-tool
	./manifest-tool push from-args --platforms $(subst $(space),$(comma),$(ALL_PLATFORMS)) --template $(DOCKER_REPO):$(VERSION)-ARCH --target $(DOCKER_REPO):$(VERSION)

push: crossbuild manifest-tool $(addprefix push-,$(ALL_ARCH)) manifest-push

.PHONY: all clean crossbuild build container push push-% manifest-push
