NAME=tomaskral/kompose-demo-frontend
VERSION=latest


.PHONY: image push


image:
	docker build -t $(NAME):$(VERSION) .

push:
	docker push $(NAME):$(VERSION)
