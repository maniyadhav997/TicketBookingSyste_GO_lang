FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

# Pure Go SQLite driver does not need CGO.
RUN CGO_ENABLED=0 go build -o ticket-system .

FROM alpine:3.19

WORKDIR /app

# Create database folder so SQLite file can be created at runtime.
RUN mkdir -p database

COPY --from=builder /app/ticket-system .

ENV PORT=8080
EXPOSE 8080

CMD ["./ticket-system"]
