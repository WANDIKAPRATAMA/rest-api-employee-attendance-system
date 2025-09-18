# Stage 1: build
FROM golang:1.24-alpine AS builder

# Atur direktori kerja di dalam container
WORKDIR /app

# Salin file-file yang diperlukan untuk go mod download
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode sumber (pastikan ada .dockerignore)
COPY . .

# Build binary dengan CGO_ENABLED=0 untuk membuat binary statis
# Ini penting agar binary bisa berjalan di Alpine tanpa dependensi glibc
RUN CGO_ENABLED=0 go build -o main ./cmd/web



# Stage 2: final minimal image
FROM alpine:3.18

# Atur direktori kerja
WORKDIR /app

# Salin binary dan config.json dari stage 'builder'
COPY --from=builder /app/main .
# COPY --from=builder /app/prod.config.json .
# Salin file konfigurasi (jika ada, misalnya .env)
# COPY --from=builder /app/.env .env

# Expose port
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]