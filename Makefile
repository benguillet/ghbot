.PHONY: build

build:
	godep go build -o build/ghbot ./cmd/ghbot
