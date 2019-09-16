FROM golang:latest as builder
LABEL maintainer="Tobias Suckow <tobias@suckow.biz>"
WORKDIR /src
COPY . /src/
RUN make

FROM alpine:latest
RUN apk --update add ca-certificates
COPY --from=builder /src/binary .
CMD ["./binary"]