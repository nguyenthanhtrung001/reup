# syntax=docker/dockerfile:1

FROM golang:1.24.1-alpine as builder
WORKDIR /app
COPY  . .
RUN go build -o consumer ./cmd/consumer

FROM alpine:latest
RUN apk add --no-cache ca-certificates && update-ca-certificates
ENV TZ=Asia/Ho_Chi_Minh
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
WORKDIR /
COPY --from=builder /app/consumer .

EXPOSE 80

ENTRYPOINT [ "./consumer" ]
