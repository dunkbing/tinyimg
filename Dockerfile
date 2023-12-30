FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN apk add --no-cache gcc musl-dev
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -o tinyimg
EXPOSE 8080
ENTRYPOINT ["./tinyimg"]
