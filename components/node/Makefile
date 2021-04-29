GO=go
GOFLAG=-a -installsuffix cgo
PROJ=node

repo=gitlab.jiangxingai.com/jxserving/components/$(PROJ)
ver=$(shell cat ./VERSION)
time=$(shell date "+%m/%d/%Y %R %Z")
hash=$(shell git rev-parse --short HEAD)
gover=$(shell go version)

LDFLAG=-ldflags '-X "$(repo)/pkg/config.buildVersion=$(ver)" -X "$(repo)/pkg/config.buildTime=$(time)" -X "$(repo)/pkg/config.buildHash=$(hash)" -X "$(repo)/pkg/config.goVersion=$(gover)"'

.PHONY: linux-amd64 linux-arm64 windows-amd64 darwin-amd64

linux-amd64: main.go
	CGO_ENABED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAG) $(LDFLAG) -o go-$@ $<
	mv go-$@ bin
	docker build -t $(PROJ)-$@:$(ver) .
	rm bin

linux-arm64: main.go
	CGO_ENABED=0 GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAG) $(LDFLAG) -o go-$@ $<
	mv go-$@ bin
	docker build -t $(PROJ)-$@:$(ver) .
	rm bin

windows-amd64: main.go
	CGO_ENABED=0 GOOS=windos GOARCH=amd64 $(GO) build $(GOFLAG) $(LDFLAG) -o go-$@ $<

darwin-amd64: main.go
	CGO_ENABED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAG) $(LDFLAG) -o go-$@ $<
