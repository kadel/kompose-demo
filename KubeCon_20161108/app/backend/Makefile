NAME=tomaskral/kompose-demo-backend
VERSION=kubecon2016


.PHONY: image push

all: build

build: server
	go build server.go

image: build
	docker build -t $(NAME):$(VERSION) .

push:
	docker push $(NAME):$(VERSION)
