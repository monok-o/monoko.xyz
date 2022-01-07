FROM golang:1.17-alpine AS builder
WORKDIR /build
RUN apk --no-cache add ca-certificates
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build main.go

# running app
FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /build /app/

ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=3000

CMD ["/app/main"]