

.PHONY: work
work:
	@go work init
	@go work use -r ./

.PHONY: build
build:
	go build -o out/cron
	chmod -R +x ./out

.PHONY: run
run:
	CONFIG=./tests/config go run main.go

.PHONY: lint
lint:
	@golangci-lint run -v
