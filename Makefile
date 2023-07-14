build:
	@go build -o bin/url-shortener

run: build
	@./bin/url-shortener

test:
	@go test -v ./...