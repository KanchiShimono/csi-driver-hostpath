REV=$(shell git describe --long --tags --match='v*' --dirty 2>/dev/null || git rev-list -n1 HEAD)

bin/%: cmd/%/main.go
	go build -ldflags '-X "main.version=$(REV)" -extldflags "-static"' -o $@ $<

.PHONY: build
build: bin/hostpathplugin
