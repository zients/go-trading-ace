FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o main .

# alpine environment
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main main
COPY --from=builder /app/config config
COPY --from=builder /app/docs/swagger.json /docs/swagger.json

CMD ["/main"]
