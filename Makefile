build:
	go build -o cmd/rest/rest_binary cmd/rest/*.go
	go build -o cmd/consumer/consumer_binary cmd/consumer/*.go

migrate:
	go run storage/migrations/main.go

run-consumer:
	./cmd/consumer/consumer_binary

run-server:
	./cmd/rest/rest_binary
