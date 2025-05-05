build:
	@templ generate views/
	@go build

run:
	@templ generate views/
	@go run main.go

templ:
	@templ generate views/

test:
	@go test -v ./...

clean:
	@rm -fr bin/
