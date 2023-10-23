FROM golang:alpine as builder
WORKDIR /app
COPY . /app
RUN go mod tidy
RUN go build -o leagueOfGoonsAPI


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/leagueOfGoonsAPI /app/
COPY --from=builder /app/scoring.toml /app/
COPY --from=builder /app/.env /app/
EXPOSE 8080
ENTRYPOINT ["/app/leagueOfGoonsAPI"]