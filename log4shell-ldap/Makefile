-include ../common.mk

NAME := log4shell-ldap

docker-build:
	docker build --platform linux/amd64 -t $(HUB)/$(NAME):$(TAG) .

docker-push:
	docker push $(HUB)/$(NAME):$(TAG)
