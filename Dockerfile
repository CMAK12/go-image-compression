FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd/app/main.go

CMD ["/app/main"]
