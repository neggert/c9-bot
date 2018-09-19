GO=go
DEP=dep
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GODEP=$(DEP) ensure

BUILD_DIR=build
BINARY_NAME=c9bot
BINARY_MAC=$(BUILD_DIR)/darwin/$(BINARY_NAME)
BINARY_LINUX=$(BUILD_DIR)/linux/$(BINARY_NAME)

VENDOR_DIR=vendor
MAIN=github.com/neggert/c9-bot/cmd/c9bot
PKGS:=$(shell go list ./...)
DOCKER_COMPOSE_FILE=deployments/docker-compose.yml

DEV_DOCKER_MACHINE=default
PROD_DOCKER_MACHINE=c9bot-prod
DEV_DOCKER_PROJECT=c9bot-dev
PROD_DOCKER_PROJECT=c9bot-prod
DEV_SECRETS=secrets/dev.env
PROD_SECRETS=secrets/prod.env

all: test build

build: $(BINARY_MAC)

$(BINARY_MAC): deps
	$(GOBUILD) -v -o $@ $(MAIN)

$(BINARY_LINUX): deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o $@ $(MAIN)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(VENDOR_DIR)

deps:
	$(GODEP)

test: deps
	$(GOTEST) $(PKGS) -v 

dev:
	eval $(bin/env-compose.sh $(DEV_DOCKER_MACHINE));\
	source $(DEV_SECRETS);\
	docker-compose -f $(DOCKER_COMPOSE_FILE) -p $(DEV_DOCKER_PROJECT) up -d 

deploy: docker docker-db
	eval $(bin/env-compose.sh $(PROD_DOCKER_MACHINE));\
	source $(PROD_SECRETS);\
	docker-compose -f $(DOCKER_COMPOSE_FILE) -p $(PROD_DOCKER_PROJECT) up -d 
	

.PHONY: clean test docker docker-db
