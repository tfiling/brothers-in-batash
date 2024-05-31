LAST_TAG?=$(shell git describe --tags 2>/dev/null || echo 'latest')
IMAGE_TAG?=$(LAST_TAG)-$(shell git branch --show-current)
LOCAL_REPO := "ldg"

export DOCKER_BUILDKIT:=1

.PHONY: build
build: build-ws

.PHONY: build-ws
build-ws:
	@echo "====================== building ws ======================"
	docker build --target webserver -t $(LOCAL_REPO)/webserver:$(IMAGE_TAG) .
	@echo "====================== building ws completed ======================"

.PHOMY: run
run: test
	@echo "====================== Running Local Dev Env ======================"
	@TAG=${IMAGE_TAG} docker compose -f deploy/local/compose.yaml up -d

.PHONY: stop
stop:
	@echo "====================== Stopping Local Dev Env ======================"
	@TAG=${IMAGE_TAG} docker compose -f deploy/local/compose.yaml down --remove-orphans -t 0

.PHONY: test
test: build
	@echo "====================== Running Tests ======================"
	docker build . --target unit-test --tag $(LOCAL_REPO)/webserver-tests:latest
	@echo "====================== Completed Running Tests ======================"
