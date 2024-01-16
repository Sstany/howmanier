FROM golang:1.21.5-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN apk update && \
    apk add build-base

RUN go mod download

COPY ./ ./

RUN go build -o howmanier /app/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/howmanier /app/

CMD [ "./howmanier" ]