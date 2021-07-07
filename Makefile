xxx:
	@echo "Please select optimal option."

build:
	@go build -o ginparam .

clean:
	@rm -f ./ginparam

run:
	@make build
	@./ginparam

test:
	@go test -v "./..."

lint:
	@go vet
