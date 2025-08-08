FROM golang:1.23-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /app/hajimi-go ./scanner

FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

USER appuser

WORKDIR /app

COPY --from=builder /app/hajimi-go ./

CMD ["/app/hajimi-go"]
