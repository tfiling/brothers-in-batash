LAST_TAG?=$(shell git describe --tags 2>/dev/null || echo 'latest')
IMAGE_TAG?=$(LAST_TAG)-$(shell git branch --show-current)
LOCAL_REPO := "ldg"

export DOCKER_BUILDKIT:=1

.PHONY: build
build: build-ws build-frontend

.PHONY: build-ws
build-ws:
	@echo "====================== building ws ======================"
	docker build -f backend/Dockerfile --target webserver -t $(LOCAL_REPO)/webserver:$(IMAGE_TAG) ./backend
	@echo "====================== building ws completed ======================"

.PHONY: build-frontend
build-frontend:
	@echo "====================== building frontend ======================"
	docker build -f frontend/Dockerfile -t $(LOCAL_REPO)/frontend:$(IMAGE_TAG) ./frontend
	@echo "====================== building frontend completed ======================"

.PHONY: run
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
	docker build -f backend/Dockerfile --target unit-test --tag $(LOCAL_REPO)/webserver-tests:latest ./backend
	@echo "====================== Running Frontend Tests ======================"
	docker run --rm -v $(PWD)/frontend:/app -w /app node:18-alpine npm test
	@echo "====================== Completed Running Tests ======================"
