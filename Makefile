GO=go

default: build

build:
	go build -o objrebuild ./cmd/objrebuild.go

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o objrebuild ./cmd/objrebuild.go