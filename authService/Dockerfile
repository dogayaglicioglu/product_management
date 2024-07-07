FROM golang:1.22.3-alpine AS builder

WORKDIR /usr/src/app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./tmp/main ./cmd/main.go

EXPOSE 8082

RUN go mod tidy

CMD ["./tmp/main"]
