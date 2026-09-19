.PHONY: test vet build check

test:
	go test ./...

vet:
	go vet ./...

build:
	go build ./cmd/localwp

check: test vet build
