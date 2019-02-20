.PHONY: all build docker-build clean

all: build

build:
	go get -d ./...
	CGO_ENABLED=0 go build ./...

docker-build:
	docker build -f Dockerfile -t hilink-reconnect ..

clean:
	rm hilink-reconnect
