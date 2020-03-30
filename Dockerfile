FROM golang:1.13-alpine3.11

RUN apk add --no-cache gcc musl-dev
RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]