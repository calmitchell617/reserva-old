include .envrc

.PHONY: vendor
vendor:
	go mod tidy
	go mod verify
	go mod vendor

.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	GOOS=linux GOARCH=arm64 go build -ldflags="-s" -o=./bin/linux_arm64/api ./cmd/api