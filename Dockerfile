FROM golang:latest as builder
LABEL maintainer="Tobias Suckow <tobias@suckow.biz>"
WORKDIR /src
COPY . /src/
RUN make

FROM scratch
COPY --from=builder /src/binary .
CMD ["./binary"]