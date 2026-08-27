FROM golang:1.26.4-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ pkg/
COPY services/media/ services/media/

RUN CGO_ENABLED=0 GOOS=linux go build -o /media-service ./services/media/cmd/api/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /media-service .

EXPOSE 8085
CMD [ "./media-service" ]