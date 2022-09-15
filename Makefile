build:
	GOOS=darwin GOARCH=arm64 go build .

install: build
	cp photo-file-filter /usr/local/bin
