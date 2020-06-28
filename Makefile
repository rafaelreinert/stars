.PHONY: build  up test prune

build:
	docker-compose build app

up:
	docker-compose up app

run: build up

prune:
	docker system prune -af

test:
	docker-compose up -d db
	go test -p 1 ./...
	docker-compose down
	