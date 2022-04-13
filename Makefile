#!make
include local_config.env

back_start:
	@docker-compose --env-file ./local_config.env up

front_start:
	cd ./web && pm2 start "npm run serve" --name web-interface

front_stop:
	cd ./web && pm2 delete web-interface

back_stop:
	@docker-compose down

server_log_show:
	@docker-compose logs -f go-docker-image

start:
	$(MAKE) front_start
	$(MAKE) back_start
	$(MAKE) server_log_show

stop:
	$(MAKE) front_stop
	$(MAKE) back_stop

db_create:
	./build/migrate-cli db:create -d postgres -u "user=${DB_USER} password=${DB_PASSWORD} host=${DB_HOST} sslmode=disable" -n ${DB_NAME}


migrate_create:
	./build/migrate-cli migrate:create -n $(NAME)

migrate_up:
	./build/migrate-cli migrate -d postgres -u ${DATABASE_URL}

migrate_down:
	./build/migrate-cli migrate:down -d postgres -u ${DATABASE_URL}

#test:
#	DATABASE_URL=testttt go test -v -race -timeout 30s ./cmd/...

#docker-compose exec go-docker-image go test ./cmd...

test:
	@docker run --name test-mysql -e MYSQL_ROOT_PASSWORD=1 -p 52000:3306 -d mysql:8.0
	@docker-compose exec go-docker-image ENV=test go test -v -race -timeout 30s ./cmd/...


#@docker stop test-mysql