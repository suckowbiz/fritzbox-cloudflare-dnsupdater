all: tests tidy binary

docker-image:
	docker-compose build --no-cache

test: tests

tests:
	go test ./...

tidy:
	go mod tidy
	go mod verify

# https://medium.com/@diogok/on-golang-static-binaries-cross-compiling-and-plugins-1aed33499671
# Build with "netgo" that enforces the use of go DNS resolver. Without the error.
binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-a \
		-tags netgo \
		-ldflags '-w -extldflags "-static"' \
		-o binary \
		.