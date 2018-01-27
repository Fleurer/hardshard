build:
	dep ensure
	go build -o ./hardshard ./main.go
test:
	go test */*_test.go
