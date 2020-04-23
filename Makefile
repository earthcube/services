BINARY := services
DOCKERVER :=`cat VERSION`
.DEFAULT_GOAL := linux
ORG := geocodes

services:
	cd cmd/$(BINARY) ; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o $(BINARY)

docker:
	docker build  --tag="$(ORG)/p418services:$(DOCKERVER)"  --file=./build/Dockerfile .

dockerlatest:
	docker build  --tag="$(ORG)/p418services:latest"  --file=./build/Dockerfile .
