FROM curlimages/curl as base
WORKDIR /app
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o yt-dlp
RUN chmod a+rx yt-dlp

FROM golang:1.22-alpine
RUN apk add --no-cache gcc musl-dev pngquant jpegoptim imagemagick libwebp-tools vips-tools yt-dlp
WORKDIR /app
COPY --from=base /app/yt-dlp /usr/local/bin
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -o tinyimg
EXPOSE 8080
ENTRYPOINT ["./tinyimg"]
