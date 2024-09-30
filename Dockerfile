FROM golang:latest AS builder

WORKDIR /app/

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o songs_lib ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app/songs_lib /app/songs_lib
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/.env /app/.env

WORKDIR /app/

CMD ["./songs_lib"]
