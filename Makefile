NAME=nubis-bastionsshkey
GO15VENDOREXPERIMENT := 1

all: build

build:
	rm -rf build && mkdir -p build
	mkdir -p build/linux && GOOS=linux go build -o build/linux/$(NAME) ./*.go
	mkdir -p build/darwin && GOOS=darwin go build -o build/darwin/$(NAME) ./*.go

clean:
	rm -rf build
