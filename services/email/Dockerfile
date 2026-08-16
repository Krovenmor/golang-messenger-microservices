FROM golang:1.26.4-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ pkg/
COPY services/email/ services/email/

RUN CGO_ENABLED=0 GOOS=linux go build -o /email-service ./services/email/cmd/api/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /email-service .

CMD [ "./email-service" ]