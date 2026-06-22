FROM golang:1.21-alpine
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o user-service ./cmd
EXPOSE 8081
CMD ["./user-service"]
