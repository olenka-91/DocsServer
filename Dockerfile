FROM golang:latest AS builder

WORKDIR /app
RUN go version

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o DocsServer ./cmd/app/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/DocsServer /app/
COPY .env /app/
EXPOSE 8080
  
CMD ["./DocsServer"]