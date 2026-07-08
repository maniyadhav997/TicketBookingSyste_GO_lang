FROM golang:1.22-alpine AS builder

WORKDIR /app

# go-sqlite3 needs CGO, so we install C build tools in Alpine.
RUN apk add --no-cache gcc musl-dev

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

# Build with CGO enabled so sqlite3 driver works.
RUN CGO_ENABLED=1 go build -o ticket-system .

FROM alpine:3.19

WORKDIR /app

# Create database folder so SQLite file can be created at runtime.
RUN mkdir -p database

COPY --from=builder /app/ticket-system .

ENV PORT=8080
EXPOSE 8080

CMD ["./ticket-system"]
