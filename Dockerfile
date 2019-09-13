FROM golang:latest
LABEL maintainer="Tobias Suckow <tobias@suckow.biz>"

WORKDIR /app
COPY . /app/
RUN make

CMD ["/app/binary"]
