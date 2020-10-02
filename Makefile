VERSION ?= $(shell git describe --match 'v[0-9]*' --tags --always)

build:
	@go build -ldflags "-X main.version=$(VERSION)"

version:
	@echo $(VERSION)

generate-docs:
	@go run ./vendor/github.com/swaggo/swag/cmd/swag/main.go init -g pkg/api/api.go

run:
	@docker start go-api-database &> /dev/null || docker run -d \
	 --name go-api-database \
	 -e POSTGRES_PASSWORD=secret \
	 -p 5432:5432 postgres &> /dev/null
	@DB_CONNECTION='postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable' \
	go run -ldflags "-X main.version=$(VERSION)" main.go

docker:
	@docker build -t go-api:$(VERSION) .

test:
	@go test -v ./...

generate:
	@go generate ./...

tools:
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
