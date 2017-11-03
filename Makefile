PROJECTNAME		=	$$(basename $$(pwd))

DOCKERHOST		=	0.0.0.0:65535
DOCKER			=	docker -H $(DOCKERHOST)
COMPOSE			=	docker-compose -H $(DOCKERHOST)
GO				?=	go
GFLAGS			?=	""
GOOS 			=	$$($(GO) env GOOS)
GOARCH			=	$$($(GO) env GOARCH)
PGKRESTORE		=	$(GO) list -f '{{range .TestImports}}{{.}} {{end}}'
EXENAME			=	$(PROJECTNAME)

SRC				=	$(wildcard model/*.go) \
					$(wildcard server/*.go) \
					$(wildcard pubsub/*.go) \
					$(wildcard sendgrid/*.go)

BIN				=	bin
EXE				=	$(BIN)/maild

default: $(EXE)


up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

$(BIN):
	mkdir -p $(BIN)

$(EXE): $(SRC) $(BIN)
	$(GO) get -d ./...
	env GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(EXE) -gcflags $(GFLAGS) ./cmd/$(PROJECTNAME)

test:
	$(PGKRESTORE) ./sendgrid | xargs $(GO) get
	$(PGKRESTORE) ./server | xargs $(GO) get
	$(GO) test -v ./...

clean:
	$(GO) clean
	rm -rf $(BIN)

image: $(EXE)
	$(DOCKER) build -t vbogretsov/maild:1 .
