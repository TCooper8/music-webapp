FROM golang

WORKDIR /usr/src/music-webapp
COPY . /usr/src/music-webapp

RUN go build
RUN go test

EXPOSE 8080

CMD ["./music-webapp"]
