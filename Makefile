.ONESHELL:
TAB=echo "\t"
CURRENT_DIR = $(shell pwd)

help:
	@$(TAB) make up - запустить сервер

up:
	go run ./cmd/gophermart/main.go -a=localhost:41839 -d=postgresql://developer:developer@localhost:5432/postgres?sslmode=disable

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)
