include common.mk

APPS := log4shell-ldap vulnerable-app wasm-patch

docker-build: $(APPS:%=docker-build/%)
docker-push: $(APPS:%=docker-push/%)

docker-build/%:
	$(MAKE) -C $(@F) $(@D)

docker-push/%:
	$(MAKE) -C $(@F) $(@D)

clean:
	$(MAKE) -C wasm-patch clean
