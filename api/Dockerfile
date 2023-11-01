FROM golang:alpine as builder
WORKDIR /app
COPY api/ /app
RUN go mod tidy
RUN go build -o leagueOfGoonsAPI


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/leagueOfGoonsAPI /app/
COPY --from=builder /app/scoring.toml /app/
COPY .env /app
EXPOSE 8080
ENTRYPOINT ["/app/leagueOfGoonsAPI"]