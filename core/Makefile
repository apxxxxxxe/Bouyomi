appname := core

sources := $(wildcard *.go)

build = GOOS=$(1) GOARCH=$(2) go build -ldflags "-H windowsgui" -o ../$(appname).exe

.PHONY: all windows

all: core_windows

core_windows: $(sources)
	$(call build,windows,amd64)
