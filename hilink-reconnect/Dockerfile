FROM golang:latest

ENV DOCKER true

WORKDIR /go/src/app
COPY hilink-reconnect .

RUN go get -d -v ./...
RUN go install -v ./...

COPY instances /etc

ENTRYPOINT ["app"]
