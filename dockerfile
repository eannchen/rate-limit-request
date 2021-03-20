FROM golang:1.15.4-alpine3.12 as builder

ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# ---------------------
FROM alpine:3.12

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
# CMD tail -f /dev/null