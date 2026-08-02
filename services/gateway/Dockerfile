FROM golang:1.26.4-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ pkg/
COPY services/gateway/ services/gateway/

RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway-service ./services/gateway/cmd/api/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /gateway-service .

EXPOSE 8080
CMD [ "./gateway-service" ]