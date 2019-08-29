all: tests tidy binary

docker-image:
	docker-compose build --no-cache

test: tests

tests:
	go test ./...

tidy:
	go mod tidy
	go mod verify

binary:
	go build -o binary .