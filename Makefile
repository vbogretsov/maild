# TODO:
# - test commands

PROJECTNAME		=	$$(basename $$(pwd))

DOCKERHOST		=	0.0.0.0:65535
DOCKER			=	docker -H $(DOCKERHOST)
COMPOSE			=	docker-compose -H $(DOCKERHOST)
GO				?=	go
PKGRESTRE		=	$(GO) get -d ./...
DBNAME			=	$(PROJECTNAME)
EXENAME			=	$(PROJECTNAME)
CONTAINERDB		=	$(DOCKER) ps --format "{{.Names}}" -f name=$(PROJECTNAME)_db
PGEXEC			=	$(DOCKER) exec $$($(CONTAINERDB))
PSQL			=	$(PGEXEC) psql -U postgres
PGREADY			=	$(PGEXEC) pg_isready
CREATEDB		=	$(PGEXEC) createdb -U postgres
DROPDB			=	$(PGEXEC) dropdb -U postgres
FINDDB			=	$(PSQL) -lt | awk '{print $$1}' | grep
PORTREADY		=	nc -z 127.0.0.1
SQLMIGRATE		=	sql-migrate

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

createdb:
	$(FINDDB) $(DBNAME) || $(CREATEDB) $(DBNAME)
	$(SQLMIGRATE) up

dropdb:
	$(FINDDB) $(DBNAME) && $(DROPDB) $(DBNAME)

$(BIN):
	mkdir -p $(BIN)

build: $(SRC) $(BIN)
	$(PKGRESTRE)
	$(GO) build -o $(BIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME)

clean:
	rm -rf $(BIN)
