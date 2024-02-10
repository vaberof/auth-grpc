PROJECT_NAME=auth-grpc

WORK_DIR_LINUX=./cmd/authgrpc
CONFIG_DIR_LINUX=./cmd/authgrpc/config

WORK_DIR_WINDOWS=.\cmd\authgrpc
CONFIG_DIR_WINDOWS=.\cmd\authgrpc\config

CURRENT_DIR=$(shell pwd)

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=admin
POSTGRES_DATABASE=auth_service

DB_URL="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable"

MIGRATIONS_PATH=migrations

run.linux: build.linux
	go run $(WORK_DIR_LINUX)/*.go \
		-config.files $(CONFIG_DIR_LINUX)/application.yaml \
		-env.vars.file $(CONFIG_DIR_LINUX)/sample.env \

build.linux: build.linux.clean
	mkdir -p $(WORK_DIR_LINUX)/build
	go build -o $(WORK_DIR_LINUX)/build/main $(WORK_DIR_LINUX)/*.go
	cp -R $(CONFIG_DIR_LINUX)/* $(WORK_DIR_LINUX)/build

build.linux.local: build.linux.clean
	mkdir -p $(WORK_DIR_LINUX)/build
	go build -o $(WORK_DIR_LINUX)/build/main $(WORK_DIR_LINUX)/*.go
	cp -R $(CONFIG_DIR_LINUX)/* $(WORK_DIR_LINUX)/build
	@echo "build.local: OK"

build.linux.clean:
	rm -rf $(WORK_DIR_LINUX)/build

run.windows:
	go run $(WORK_DIR_WINDOWS)\. \
		-config.files $(CONFIG_DIR_WINDOWS)\application.yaml \
		-env.vars.file $(CONFIG_DIR_WINDOWS)\sample.env

migrate.up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate.down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

proto.gen:
	rm -rf genproto
	./scripts/gen_proto.sh ${CURRENT_DIR}
