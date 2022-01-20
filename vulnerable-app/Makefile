-include ../common.mk

NAME := vulnerable

docker-build:
	docker build --platform linux/amd64 -t $(HUB)/$(NAME):$(TAG) .

docker-push:
	docker push $(HUB)/$(NAME):$(TAG)
