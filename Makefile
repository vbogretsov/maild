# TODO:
# - test commands

PROJECTNAME		=	$$(basename $$(pwd))

DOCKERHOST		=	0.0.0.0:65535
DOCKER			=	docker -H $(DOCKERHOST)
COMPOSE			=	docker-compose -H $(DOCKERHOST)
GO				?=	go
GFLAGS			=	"-N -l"
PKGRESTRE		=	$(GO) get -d ./...
DBNAME			=	$(PROJECTNAME)
EXENAME			=	$(PROJECTNAME)
CONTAINERDB		=	$(DOCKER) ps --format "{{.Names}}" -f name=$(PROJECTNAME)_db
DBUSER			=	postgres
PGEXEC			=	$(DOCKER) exec $$($(CONTAINERDB))
PSQL			=	$(PGEXEC) psql -U $(DBUSER)
PGREADY			=	$(PGEXEC) pg_isready
CREATEDB		=	$(PGEXEC) createdb -U $(DBUSER)
DROPDB			=	$(PGEXEC) dropdb -U $(DBUSER)
FINDDB			=	$(PSQL) -lt | awk '{print $$1}' | grep
PORTREADY		=	nc -z 127.0.0.1
# go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
# go build -o bin/migrate -tags 'postgres' github.com/mattes/migrate/cli
SQLMIGRATE		=	bin/dbmigrate

SRC				=	$(wildcard model/*.go) $(wildcard server/*.go)
BIN				= 	bin

define wait
	@while $1 ; do sleep 1; done
endef

define wait_port
	$(call wait,! $(PORTREADY) $1)
endef

define wait_pg
	$(call wait,! $(PGREADY))
	@sleep 1
	$(call wait,! $(PGREADY))
endef

default: build

up:
	$(COMPOSE) up -d
	$(call wait_pg)

down:
	$(COMPOSE) down

migratetool: $(BIN)
	$(GO) get -d github.com/mattes/migrate/cli github.com/lib/pq
	$(GO) build -o $(SQLMIGRATE) -tags 'postgres' github.com/mattes/migrate/cli

createdb: migratetool
	$(FINDDB) $(DBNAME) || $(CREATEDB) $(DBNAME)
	$(SQLMIGRATE) -source file://db/migrations -database postgres://$(DBUSER)@localhost/$(DBNAME)?sslmode=disable up

dropdb:
	$(FINDDB) $(DBNAME) && $(DROPDB) $(DBNAME)

$(BIN):
	mkdir -p $(BIN)

build: $(SRC) $(BIN)
	$(PKGRESTRE)
	$(GO) build -o $(BIN)/$(PROJECTNAME) -gcflags $(GFLAGS) ./cmd/$(PROJECTNAME)

clean:
	rm -rf $(BIN)
