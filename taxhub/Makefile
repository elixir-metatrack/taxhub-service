GOPATH := $(shell go env GOPATH)

.PHONY: build run swag docker test integration clean

build:
	go build -o bin/taxhub ./cmd/taxhub

run: swag
	go run ./cmd/taxhub

swag: $(GOPATH)/bin/swag
	$(GOPATH)/bin/swag init -g cmd/taxhub/main.go
	go run ./tools/gendocs

docker:
	docker build -t taxhub:local .

test: swag
	go test ./...

integration:
	CSV_PATH=scientific_names.csv go test ./test/integration/... -tags integration -v

clean:
	rm -rf bin/ docs/

$(GOPATH)/bin/swag:
	go install github.com/swaggo/swag/cmd/swag@latest