FROM golang:1.18-alpine

WORKDIR /redis_docker

COPY . .

RUN go get github.com/go-redis/redis

RUN go build -o main

CMD [ "./main.go" ]