FROM golang:1.26.4-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ pkg/
COPY services/auth/ services/auth/

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./services/auth/cmd/api/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /auth-service .

EXPOSE 8081
CMD [ "./auth-service" ]