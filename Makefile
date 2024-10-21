make build:
	go build -race -v cmd/main.go

make test:
	go test ./... -v

make docker-build:
	docker build --tag ggr:dev .
.PHONY: docker-build

make run:
	go run -race cmd/main.go

lint:
	golangci-lint run
.PHONY: lint
