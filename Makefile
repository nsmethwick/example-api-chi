include .env.local

.DEFAULT_GOAL := run

PROJECT_APP_NAME ?= example-api

# hl is a program (see Brewfile) that parses JSON structured logs.
hl := hl -P -F -e -t "%Y.%m.%d %H:%M:%S.%3N" --flatten never

.PHONY: help
help:
	@echo "Targets:"
	@MAX_LEN_TARGET=$$(grep -h -E '^[0-9a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "; max=0} {len=length($$1); if (len>max) max=len} END {print max}'); \
	grep -h -E '^[0-9a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk -v max_len="$$MAX_LEN_TARGET" 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-" max_len "s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Variables:"
	@MAX_LEN_VAR=$$(grep -h -E '^[0-9a-zA-Z_.-]+\s[?:]?=.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "[?:]?=.*?## "; max=0} {len=length($$1); if (len>max) max=len} END {print max}'); \
	grep -h -E '^[0-9a-zA-Z_.-]+\s[?:]?=.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk -v max_len="$$MAX_LEN_VAR" 'BEGIN {FS = "[?:]?=.*?## "}; {printf "  \033[36m%-" max_len "s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Default target:"
	@printf "  \033[36m%s\033[0m\n" $(.DEFAULT_GOAL)

.PHONY: check-container-runtime
check-container-runtime:  ## Check the container runtime (docker, colima) is running
	@colima status > /dev/null 2>&1 || (echo "starting the colima container runtime" && $(MAKE) docker-start)

.PHONY: mysql-connect
mysql-connect: ## Connect to the dockerized MySQL
	@echo "Once logged in, try this command:"
	@echo "use example-api; select * from customers;"
	docker exec -it example-api-mysql mysql -u root -p example-api
	@# host version:
	@# mysql -h 127.0.0.1 -P 3306 -u root -p example-api

.PHONY: mysql-drop-db
mysql-drop-db: ## Drop all tables in the database
	@echo "Dropping database tables..."; \
	PROJECT_APP_NAME=$(strip $(PROJECT_APP_NAME)) docker-compose --project-directory local/mysql exec -T mysql mysql -u$(MYSQL_USER) -p$(MYSQL_PASSWORD) -D$(MYSQL_DATABASE) -e "DROP TABLE IF EXISTS customers;" || exit 1; \
	echo "Database tables dropped"

.PHONY: mysql-seed
mysql-seed: ## Seed the database with SQL files from initdb.d
	@echo "Seeding database..."; \
	for sqlfile in local/mysql/initdb.d/*.sql; do \
		echo "Executing $$sqlfile..."; \
		PROJECT_APP_NAME=$(strip $(PROJECT_APP_NAME)) docker-compose --project-directory local/mysql exec -T mysql mysql -u$(MYSQL_USER) -p$(MYSQL_PASSWORD) -D$(MYSQL_DATABASE) < $$sqlfile || exit 1; \
	done; \
	echo "Database seeded successfully"

.PHONY: mysql-start
mysql-start: check-container-runtime ## Spin-up mysql stack locally
	@echo "Starting MySQL..."; \
	mkdir -p local/mysql/initdb.d; \
	PROJECT_APP_NAME=$(strip $(PROJECT_APP_NAME)) docker-compose --project-directory local/mysql up --detach; \
	echo "⏳ Waiting for MySQL to be ready..."; \
	while ! PROJECT_APP_NAME=$(strip $(PROJECT_APP_NAME)) docker-compose --project-directory local/mysql exec -T mysql mysql -u$(MYSQL_USER) -p$(MYSQL_PASSWORD) -D$(MYSQL_DATABASE) -e "SELECT 1" >/dev/null 2>&1; do \
		sleep 1; \
	done; \
	echo "MySQL is ready."

.PHONY: mysql-stop
mysql-stop: check-container-runtime ## Spin-down mysql stack locally
	@echo "Stopping MySQL..."; \
	PROJECT_APP_NAME=$(strip $(PROJECT_APP_NAME)) docker-compose --project-directory local/mysql down --remove-orphans --volumes; \
	echo "MySQL stopped"

define RUN_COMMAND
$(MAKE) mysql-start
$(MAKE) mysql-drop-db;
$(MAKE) mysql-seed;
go run ./cmd/api/main.go;
$(MAKE) mysql-stop;
endef
.PHONY: run-with-deps
run-with-deps: ## Starts the app along with other dependancies
	$(RUN_COMMAND)

.PHONY: run
run: ## Clean run without automatically starting other deps or seeding db
	go run ./cmd/api/main.go