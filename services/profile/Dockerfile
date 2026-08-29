FROM golang:1.26.4-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ pkg/
COPY services/profile/ services/profile/

RUN CGO_ENABLED=0 GOOS=linux go build -o /profile-service ./services/profile/cmd/api/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /profile-service .

EXPOSE 8086
CMD [ "./profile-service" ]