-include ../common.mk

NAME := vulnerable

docker-build:
	docker build -t $(HUB)/$(NAME):$(TAG) .

docker-push:
	docker push $(HUB)/$(NAME):$(TAG)
