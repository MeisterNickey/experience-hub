.PHONY: build run-api run-worker test docker-up docker-down tf-init tf-validate

build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker

run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

test:
	go test -v -race ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down

tf-init:
	cd terraform && terraform init

tf-validate:
	cd terraform && terraform fmt -check && terraform validate
