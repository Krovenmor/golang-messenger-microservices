FROM golang:1.26.4-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ pkg/
COPY services/ws/ services/ws/

RUN CGO_ENABLED=0 GOOS=linux go build -o /ws-service ./services/ws/cmd/api/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /ws-service .

EXPOSE 8083
CMD [ "./ws-service" ]