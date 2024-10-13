# default value
steps=1
forceMigration=false

.PHONY: help
help: ## Show help command
	@printf "Makefile Command\n";
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: migrate
migrate: ## Create Migrations file, example : make migrate name="xxxx"
	@if [ -z "${name}" ]; then \
		echo "Error: name is required \t example : make migrate name="name_file_migration";" \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations '${name}'


migrate-up: ## Up migration, example : make migrate-up steps=1 forceMigration=true
	go run migrations/main.go -steps=${steps} -forceMigration=${forceMigration}

migrate-rollback: ## Rollback, example : make migrate-rollback steps=1 forceMigration=true
	go run migrations/main.go -rollback -steps=${steps}  -forceMigration=${forceMigration}

migrate-force: ## Force migration when dirty state happen, example : make migrate-force version=1
	go run migrations/main.go -force -version=${version}


build:	
	docker build . -t checkout_svc

run:
	docker run --name checkout_svc -v ./.env:/app/.env -p 8080:8080 --restart=on-failure:5 --network pg_checkout_svc_network checkout_svc

up:
	docker build . -t checkout_svc
	docker run  --add-host=host.docker.internal:host-gateway  --name checkout_svc -v ./.env:/app/.env -dp 8080:8080 --restart=on-failure:5 --network pg_checkout_svc_network checkout_svc

down:
	docker stop checkout_svc || true
	docker rm checkout_svc || true
	docker rmi checkout_svc

app-enter-image:
	docker run -it --rm --entrypoint sh checkout_svc

app-enter-container:
	docker exec -it checkout_svc sh


test-integration:
	go test -v ./test/integration/... -count=1