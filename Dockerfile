FROM golang:latest as builder
WORKDIR /src
COPY . /src/
RUN make

FROM alpine as alpine
RUN apk update \
 && apk add --no-cache git ca-certificates \
 && update-ca-certificates

FROM scratch
LABEL maintainer="Tobias Suckow <tobias@suckow.biz>"
ENV FCD_PORT=80
EXPOSE 80
COPY --from=builder /src/binary .
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["./binary"]