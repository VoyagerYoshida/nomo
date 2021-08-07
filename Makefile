IMAGE_NAME := voyagerwy130/pitcher-nomo:1.0
CONTAINER_NAME := nomo
VOLUMEDIR := /var/www/
PWD := $(shell pwd)

.PHONY: build
build:
	@docker build . -t $(IMAGE_NAME)

.PHONY: run
run:
	@docker run \
			-it --rm \
			--name=$(CONTAINER_NAME) \
			-v $(PWD)/N-21:$(VOLUMEDIR)/monitored \
			-v $(PWD)/pdf:$(VOLUMEDIR)/out \
			--env-file $(PWD)/.env \
			$(IMAGE_NAME) 

.PHONY: stop
stop:
	@docker rm -f $(CONTAINER_NAME)

.PHONY: test
test:
	@cp $(PWD)/assets/TestImage.pdf $(PWD)/N-21/TestImage.pdf
	@docker run \
			-it --rm \
			--name=$(CONTAINER_NAME) \
			-v $(PWD)/N-21:$(VOLUMEDIR)/monitored \
			-v $(PWD)/pdf:$(VOLUMEDIR)/out \
			--env-file $(PWD)/.env \
			--entrypoint="go" \
			$(IMAGE_NAME) \
			test -v ./bot
