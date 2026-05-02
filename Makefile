VERSION := 0.1.0
BINARY  := aict
MODULE  := github.com/aict-tool/aict

.PHONY: all build-all build-windows build-macos build-linux portable-windows test clean

all: build-all

test:
	go test ./...

build-all: build-windows build-macos build-linux

build-windows:
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY)-windows-amd64.exe ./cmd/aict/

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY)-macos-amd64 ./cmd/aict/
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY)-macos-arm64 ./cmd/aict/

build-linux:
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY)-linux-amd64 ./cmd/aict/

portable-windows:
	mkdir -p dist/portable-windows
	cp dist/$(BINARY)-windows-amd64.exe dist/portable-windows/$(BINARY).exe
	mkdir -p dist/portable-windows/.aict
	cd dist && zip -r $(BINARY)-portable-windows-amd64.zip portable-windows/

clean:
	rm -rf dist/
