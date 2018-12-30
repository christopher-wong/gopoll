.PHONY: all server client

server:
	go run server/main.go

client:
	go run client/main.go
