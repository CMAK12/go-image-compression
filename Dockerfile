FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd/app/main.go
EXPOSE 8080

CMD ["/app/main"]
