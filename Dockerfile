# syntax=docker/dockerfile:1

FROM golang:alpine AS builder
WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stresstest ./cmd/stresstest

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/stresstest ./

ENTRYPOINT ["./stresstest"]
