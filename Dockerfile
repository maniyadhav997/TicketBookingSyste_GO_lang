FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o ticket-system .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/ticket-system .
COPY --from=builder /app/database ./database

ENV PORT=8080
EXPOSE 8080

CMD ["./ticket-system"]

