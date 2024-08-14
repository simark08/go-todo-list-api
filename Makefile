include .env

build:
	@go build -o bin/main cmd/main.go 

run: build
	@./bin/main

up:
	cd db/migrations && goose ${DB_ENGINE} "${DB_USER}:${DB_PASSWORD}@/${DB_DATABASE}?parseTime=true" up

down:
	cd db/migrations && goose ${DB_ENGINE} "${DB_USER}:${DB_PASSWORD}@/${DB_DATABASE}?parseTime=true" down


	