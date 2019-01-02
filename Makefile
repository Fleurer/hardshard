build:
	go build -o bin/hardshard ./cmd/hardshard/main.go
test:
	go test -v -cover ./...
