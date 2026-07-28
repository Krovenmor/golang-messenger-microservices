FROM golang:1.26.4-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ pkg/
COPY services/msg/ services/msg/

RUN CGO_ENABLED=0 GOOS=linux go build -o /msg-service ./services/msg/cmd/api/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /msg-service .

EXPOSE 8082
CMD [ "./msg-service" ]