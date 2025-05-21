FROM golang:1.24 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /tcp-trace main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /tcp-trace /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/tcp-trace"]

