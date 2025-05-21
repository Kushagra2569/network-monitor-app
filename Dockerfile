FROM golang:1.23-alpine

WORKDIR /app
COPY . .

RUN go build -o fetcher main.go

CMD ["./fetcher"]
