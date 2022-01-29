create-migration:
	docker-compose run migration dbmate --migrations-dir /db/migrations new $(name)
	sudo chown -R $(USER) ./db

up:
	docker-compose down -v
	docker-compose run migration
	docker-compose up api

test:
	docker-compose down -v
	docker-compose up -d postgres migration
	docker-compose run api go test  -failfast ./tests

