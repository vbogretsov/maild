PROJECTNAME		=	$$(basename $$(pwd))
GO				?=	go
DOCKERHOST		=	0.0.0.0:65535
DOCKER			=	docker -H $(DOCKERHOST)
DOCKERCOMPOSE	=	docker-compose -H $(DOCKERHOST)
BINDIR			=	bin


all: build


$(BINDIR):
	mkdir -p $(BINDIR)

deps:
	$(GO) get -d ./...


maild: $(BINDIR) deps
	$(GO) build -o $(BINDIR)/maild ./cmd/maild

mailcd: $(BINDIR) deps
	$(GO) build -o $(BINDIR)/mailcd ./cmd/mailcd


build: maild mailcd



clean:
	$(GO) clean
	rm -rf $(BINDIR)


up:
	$(DOCKERCOMPOSE) up -d


down:
	$(DOCKERCOMPOSE) down
