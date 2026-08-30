FROM golang:alpine AS builder
WORKDIR /app
COPY scoring.toml api/ /app/
RUN go mod tidy
RUN go build -o leagueOfGoonsAPI


FROM alpine:latest
ARG ENV_FILE
WORKDIR /app
COPY --from=builder /app/leagueOfGoonsAPI /app/
COPY --from=builder /app/scoring.toml /app/
COPY ${ENV_FILE} /app/.env
EXPOSE 8080
ENTRYPOINT ["/app/leagueOfGoonsAPI"]
