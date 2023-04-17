.PHONY: proto serve migration db

proto:
	protoc -I ./proto \
	--go_out=./pb --go_opt=paths=source_relative \
	--go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
	./proto/backend.proto

serve:
	go run cmd/main.go

migration:
	go run cmd/main.go migration

db:
	docker compose down
	docker compose up -d
