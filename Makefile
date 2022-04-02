all: build

build:
	go build

install: build
	install tailscale_hosts ~/bin/
	rm tailscale_hosts
