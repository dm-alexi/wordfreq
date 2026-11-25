.PHONY: check

check:
	go test ./... -race
	go vet ./...
	staticcheck ./...
	gofmt -l .
