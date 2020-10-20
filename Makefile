SOURCES = $(shell find . -name "*.go")

default: build

build:
	go build ./...

check:
	go test -count 1 ./...

imports:
	@goimports -w $(SOURCES)

fmt:
	@gofmt -w -s $(SOURCES)

.coverprofile: $(SOURCES)
	go test -count 1 -coverprofile .coverprofile

cover: .coverprofile
	go tool cover -func .coverprofile

showcover: .coverprofile
	go tool cover -html .coverprofile
