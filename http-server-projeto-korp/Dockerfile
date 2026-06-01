FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /http-server-projeto-korp .

FROM alpine:3.20

RUN adduser -D -g '' appuser

WORKDIR /app
COPY --from=build /http-server-projeto-korp .

USER appuser

EXPOSE 8080

ENTRYPOINT ["./http-server-projeto-korp"]
