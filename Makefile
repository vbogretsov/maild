PROJECTNAME		=	$$(basename $$(pwd))
GO				?=	go
DOCKERHOST		=	0.0.0.0:65535
DOCKER			=	docker -H $(DOCKERHOST)
DOCKERCOMPOSE	=	docker-compose -H $(DOCKERHOST)
BINDIR			=	bin


all: build


$(BINDIR):
	mkdir -p $(BINDIR)


build: $(BINDIR)
	$(GO) get -d ./...
	$(GO) build -o $(BINDIR)/maild ./cmd/maild


clean:
	$(GO) clean
	rm -rf $(BINDIR)


up:
	$(DOCKERCOMPOSE) up -d


down:
	$(DOCKERCOMPOSE) down
