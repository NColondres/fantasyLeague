FROM golang:alpine AS builder
WORKDIR /app
COPY scoring.toml sentinel/ /app/
RUN go mod tidy
RUN go build -o sentinel


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/sentinel /app/
COPY --from=builder /app/scoring.toml /app/
COPY .env /app
ENTRYPOINT ["/app/sentinel"]