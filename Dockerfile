FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download 

COPY . .
COPY tls ./tls

RUN go build -o snippetbox ./cmd/web

EXPOSE 4000

CMD ["./snippetbox"]