include .env
export

service-run:
	go run ./cmd/main.go

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down