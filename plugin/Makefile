STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo

build:
	GOPRIVATE="github.com/opengovern" CC=/usr/bin/musl-gcc GOPRIVATE="github.com/opengovern" GOOS=linux GOARCH=amd64 go build -v -ldflags "-linkmode external -extldflags '-static' -s -w" -o ./build/cohereai.plugin *.go


install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/turbot/cohereai@latest/cohereai.plugin -tags "${BUILD_TAGS}" *.go
