-include ../Makefile.mk

NAME := log4shell-ldap

docker-build:
	docker build -t $(HUB)/$(NAME):$(TAG) .

docker-push:
	docker push $(HUB)/$(NAME):$(TAG)
