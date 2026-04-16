# Stage 1: Build
FROM golang:alpine AS builder

# Pasang git karena beberapa modul Go membutuhkannya untuk download
RUN apk add --no-cache git

WORKDIR /app

# Menggunakan wildcard (*) agar jika go.sum tidak ada, build tidak gagal
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# Build dengan flags agar binary lebih ringan
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Final Image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/main .

EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]