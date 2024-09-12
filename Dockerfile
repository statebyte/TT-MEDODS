FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o auth_service ./src/cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/auth_service .

COPY --from=build /app/migrations ./migrations

RUN apk --no-cache add bash

CMD ["./auth_service"]
