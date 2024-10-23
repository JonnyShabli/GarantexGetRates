build:
	docker build --tag ggr:dev .
.PHONY: build

test:
	go test ./... -v
.PHONY: test

docker-build:
	docker build --tag ggr:dev .
.PHONY: docker-build

run:
	go run -race cmd/main.go
.PHONY: run

lint:
	golangci-lint run
.PHONY: lint
