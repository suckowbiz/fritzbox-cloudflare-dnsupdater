FROM golang:latest

WORKDIR /app
ADD . /app/
RUN go test ./...
RUN go build -o binary .

CMD ["/app/binary"]