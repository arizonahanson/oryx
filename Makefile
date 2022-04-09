NAME=oryx
PKGPATH=github.com/arizonahanson/oryx
GITHASH=$(shell git rev-parse --short=12 HEAD)
UTCTIME=$(shell date -u '+%Y%m%d%H%M%S')
LDFLAGS=-X $(PKGPATH)/cmd.version=$(UTCTIME)-$(GITHASH)
GONAME=$(HOME)/go/bin/$(NAME)
GODEPS=*/*.go */**/*.go
PIGEON=pigeon
PEGIN=internal/parser/parser.peg
PEGOUT=internal/parser/parser.go

default: install

$(NAME): $(GODEPS) Makefile
	@echo "Building..."
	go build -ldflags="$(LDFLAGS)"

$(GONAME): $(GODEPS) $(PEGOUT) Makefile
	@echo "Compiling and installing..."
	go install -ldflags="$(LDFLAGS)"

$(PEGOUT): $(PEGIN)
	@echo "Generating parser..."
	$(PIGEON) -o "$(PEGOUT)" "$(PEGIN)"

.PHONY: build
build: $(NAME)

.PHONY: install
install: $(GONAME)
