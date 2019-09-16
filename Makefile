all: tests tidy binary

docker-image:
	docker-compose build --no-cache

test: tests

tests:
	go test ./...

tidy:
	go mod tidy
	go mod verify

# Build with:
# - a  				to force build
# - ldflags '-w' 	do not include debug information to keep file size low
# - "netgo" 		enforces the use of go DNS resolver and resolves "standard_init_linux.go:211: exec user process caused "no such file or directory"
binary:
	GOOS=linux GOARCH=amd64 go build \
		-a \
		-tags netgo \
		-ldflags '-w' \
		-o binary \
		.