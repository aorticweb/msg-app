# msg-app
coding challenge

## Context
I tried to keep the implementation light, use the default golang packages,
there are multiple areas of improvement listed further down this page.
All the tests are integration tests since the Postman collections provides some 
form of e2e testing.
NB: this would have been quicker in python but less fun.

## requirements
- docker
- docker-compose
- psql (optional, if you want to explore the db)

## services
all services run inside containers
postgres: Postgres database
migration: runs migration for postgres database
api: go http api for message app (exposes 3001)

note: there are two target for the api image, dev which has live code reload and prod which only contains
the api server executable

# run 
make up

# test
make test

# Area of Improvements
- some code could be abstracted (lot of DRY violations)
- add more integration testing (not every scenario is covered)
- add unit testing of crud, model and handlers
- add e2e tests using a different language (JS or python)
- add created, updated columns for audit log
- db does not allow empty messages
- soft delete message
- switch from id int autoincrement to uuid
- Add Pagination to messages GET endpoints
- improve request payload validation with to return errors with more context
