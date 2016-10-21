NAME=tomaskral/kompose-demo-backend
VERSION=latest


.PHONY: image push

all: build

build: server
	go build server.go

image: build
	docker build -t $(NAME):$(VERSION) .

push:
	docker push $(NAME):$(VERSION)
