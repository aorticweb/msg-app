create-migration:
	docker-compose run migration dbmate --migrations-dir /db/migrations new $(name)
	sudo chown -R $(USER) ./db

init-db:
	docker-compose run migration
	docker-compose up api
