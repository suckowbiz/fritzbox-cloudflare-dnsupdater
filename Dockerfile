FROM golang:latest
LABEL maintainer="Tobias Suckow <tobias@suckow.biz>"

WORKDIR /app
ADD . /app/
RUN make

CMD ["/app/binary"]