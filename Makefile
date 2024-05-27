SHELL=/usr/bin/env bash

DOCKER_COMPOSE := $(shell command -v docker-compose 2> /dev/null)
IMAGE_TAG := $(shell basename `git rev-parse --show-toplevel`)
DEBUG=false


ensure_devcontainer:
	@ mkdir -vp .devcontainer/_cache


lint:
	gofumpt -w ./..
	golangci-lint run --fix


lint_docker: ensure_devcontainer
	@ docker run -it \
		--rm \
		-v $(shell pwd):/app:ro \
		-v $(shell pwd)/.devcontainer/_cache:/go \
		-w /app \
		--network host \
		golangci/golangci-lint:v1.56.0 \
		/bin/bash -c 'git config --global --add safe.directory /app && golangci-lint --version && go get ./... && golangci-lint run -v'


golangci_docker_shell: ensure_devcontainer
	@ docker run -it \
		--rm \
		-v $(shell pwd):/app:ro \
		-v $(shell pwd)/.devcontainer/_cache:/go \
		-w /app \
		--network host \
		golangci/golangci-lint:v1.56.0 \
		/bin/bash -c 'git config --global --add safe.directory /app && bash'


golangci_docker_fix: ensure_devcontainer
	@ docker run -it \
		--rm \
		-v $(shell pwd):/app \
		-v $(shell pwd)/.devcontainer/_cache:/go \
		-w /app \
		--network host \
		golangci/golangci-lint:v1.56.0 \
		/bin/bash -c 'git config --global --add safe.directory /app && go get ./... && golangci-lint run --fix -v'
	@ make fix_permissions


test_docker: ensure_devcontainer
	@ docker run -it \
		--rm \
		-v $(shell pwd):/app:ro \
		-v $(shell pwd)/.devcontainer/_cache:/go \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-w /app \
		--privileged \
		--network host \
		golang:1.21.8-bullseye \
		/bin/bash -c 'git config --global --add safe.directory /app && go get ./... && go test -v ./...'


fix_permissions:
	@ sudo chown -vR $${USER}:$${USER} ./ | grep -v retained


build:
	@ DOCKER_BUILDKIT=1 docker build \
		--tag $(IMAGE_TAG) \
		--build-arg DEBUG=$(DEBUG) \
		--network=host \
		--output type=local,dest=build/package/gosdelnet \
		--target=export \
		./


mocker_docker:
	@ docker run -it \
		--rm \
		-v $(shell pwd):/src \
		-w /src \
		--entrypoint /bin/ash \
		vektra/mockery:v2.36
