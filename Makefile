# Making a docker container command
// Usage: make <command>
// Commands:
//   build: build the docker container
//   run: run the docker container
//   stop: stop the docker container
//   clean: remove the docker container
//   help: show this help message
//   shell: run a shell in the docker container
//   test: run the tests in the docker container

.PHONY: build run stop clean help shell test

# Variables for the docker container
IMAGE_NAME = acsp-backend-app
CONTAINER_NAME = acsp
PORT = 5000
REGISTRY_NAME = acsp


build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -d --name $(CONTAINER_NAME) -p $(PORT):$(PORT) $(IMAGE_NAME)

stop:
	docker stop $(CONTAINER_NAME)

clean:
	docker rm $(CONTAINER_NAME)

help:
@echo "Usage: make <command>"
	@echo "Commands:"
	@echo "  build: build the docker container"
	@echo "  run: run the docker container"
	@echo "  stop: stop the docker container"
	@echo "  clean: remove the docker container"
	@echo "  help: show this help message"
	@echo "  shell: run a shell in the docker container"
	@echo "  test: run the tests in the docker container"

shell:
	docker exec -it $(CONTAINER_NAME) /bin/bash

test:
	docker exec -it $(CONTAINER_NAME) /bin/bash -c "cd /app && python -m unittest discover -s tests -p '*_test.py'"

generate:
	swag init -g cmd/main.go
	docker-compose up
push:
	docker tag $(IMAGE_NAME) registry.digitalocean.com/$(REGISTRY_NAME)/$(IMAGE_NAME)
	docker push registry.digitalocean.com/$(REGISTRY_NAME)/$(IMAGE_NAME)