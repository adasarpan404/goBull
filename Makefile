# Makefile - helper tasks for development


.PHONY: test lint fmt vet ci compose-up compose-down

test:
	go test ./... -v

fmt:
	gofmt -w .
	goimports -w . || true

vet:
	go vet ./...

lint:
	golangci-lint run

compose-up:
	docker-compose -f docker-compose.ci.yml up -d --remove-orphans

compose-down:
	docker-compose -f docker-compose.ci.yml down --remove-orphans

ci: fmt vet lint compose-up
	@echo "Waiting for Redis to become healthy..."
	./scripts/wait-for-redis.sh localhost 6379
	REDIS_ADDR=localhost:6379 go test ./... -v
	$(MAKE) compose-down
