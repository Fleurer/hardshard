build:
	dep ensure
	go build -o ./hardshard ./main.go
test:
	go test -v -cover ./mysql
