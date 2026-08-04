FROM golang:1.26.4-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ pkg/
COPY services/status/ services/status/

RUN CGO_ENABLED=0 GOOS=linux go build -o /status-service ./services/status/cmd/api/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /status-service .

EXPOSE 8084
CMD [ "./status-service" ]