## Deskripsi Proyek

Proyek ini adalah backend untuk sistem absensi karyawan di perusahaan multinasional dengan >50 karyawan dan berbagai departemen. Sistem ini mencatat kehadiran, mengevaluasi kedisiplinan (ketepatan waktu berdasarkan max clock in/out per departemen), dan mendukung autentikasi.

### Requirement dari Soal Letify ID

- **ERD dan Flowchart**: Diimplementasikan dengan tabel Employee (integrasi ke UserProfile), Department, Attendance, dan AttendanceHistory.
- **Endpoint**:
  - CRUD Karyawan: Terintegrasi dengan auth (create/update/delete via UserProfile).
  - CRUD Departemen: Lengkap dengan validasi admin-only.
  - POST Absen Masuk: Clock in dengan history log.
  - PUT Absen Keluar: Clock out dengan history log.
  - GET List Log Absensi: Dengan filter tanggal/departemen, dan perhitungan ketepatan (On Time/Late/Early Leave).
- Tambahan Nilai: Integrasi autentikasi (JWT, refresh token), pagination, dynamic filters, dan dashboard metrics (e.g., total employee per dept, registrations today).

Proyek ini dibuat untuk test kerja Letify ID, memenuhi semua spesifikasi dengan arsitektur clean (layered) dan scalable.

**Repository GitHub**: Buat repo baru (e.g., `github.com/username/absensi-backend`). Push kode lengkap ke sana. Gunakan `.gitignore` untuk ignore `env`, `bin`, dll.

## Arsitektur Aplikasi

Aplikasi menggunakan arsitektur **Layered Clean Architecture** dengan pemisahan concern untuk maintainability dan testability. Struktur folder utama:

```
project-root/
├── cmd/                # Entry point (main.go)
├── internal/           # Core app logic
│   ├── app_config/     # Config app (e.g., init DB, Fiber, Logger)
│   ├── controller/     # Handler HTTP (Fiber handlers)
│   ├── dto/            # Data Transfer Objects (request/response structs)
│   ├── middleware/     # Auth middleware (JWT validation)
│   ├── models/         # Domain models (GORM structs: User, UserProfile, etc.)
│   ├── repository/     # Database interactions (GORM queries)
│   ├── route/          # API routes (Fiber groups)
│   ├── usecase/        # Business logic (services)
│   └── utils/          # Helpers (e.g., validation, response wrappers, JWT)
├── go.mod              # Go modules
├── go.sum
├── .env                # Environment variables (DB URL, JWT secret)
└── README.md           # Dokumentasi ini
```

### Layer-Layer Utama

1. **DTO (Data Transfer Objects)**:

   - Berlokasi di `internal/dto/`.
   - Berisi structs untuk request/response API (e.g., `SignupRequest`, `CreateDepartmentRequest`, `AttendanceLogResponse`).
   - Gunakan `json` tags untuk marshalling dan `validate` tags untuk validasi (go-playground/validator).
   - Tujuan: Isolasi data transfer dari domain models, cegah exposure field sensitif (e.g., password).

2. **Controller**:

   - Berlokasi di `internal/controller/`.
   - Handle HTTP requests (Fiber Ctx), bind/validate DTO, panggil UseCase, dan return response.
   - Contoh: `AuthController` untuk signup/signin, `AttendanceController` untuk clock in/out.
   - Gunakan `utils.SuccessResponse` dan `ErrorResponse` untuk standar API response (status, message, payload, errors).

3. **UseCase (Business Logic)**:

   - Berlokasi di `internal/usecase/`.
   - Layer service yang mengandung logika bisnis (e.g., validasi tambahan, transaction, agregasi).
   - Tidak bergantung pada HTTP/DB; panggil Repository untuk data.
   - Contoh: `AuthUseCase` untuk signup (hash password, create user+profile+security+role), `AttendanceUseCase` untuk clock in (cek existing, create attendance+history).

4. **Repository**:

   - Berlokasi di `internal/repository/`.
   - Interaksi langsung dengan DB via GORM (queries, CRUD, joins).
   - Contoh: `UserRepository` untuk find/create user, `AttendanceRepository` untuk find/create attendance.
   - Gunakan transactions untuk operasi multi-table (e.g., clock in + history).

5. **Models (Domain)**:

   - Berlokasi di `internal/models/`.
   - GORM structs untuk tabel (e.g., `User`, `UserProfile`, `Department`, `Attendance`).
   - Sertakan relationships (e.g., `gorm:"foreignKey"`), soft deletes (`DeletedAt`), dan default values.

6. **Middleware**:

   - Berlokasi di `internal/middleware/`.
   - Auth middleware untuk validate JWT, extract userID/role ke context.

7. **Utils**:
   - Berlokasi di `internal/utils/`.
   - Helpers seperti JWT generation, response wrappers, field generators (e.g., EmployeeCode), error details.

### Module-Module

Aplikasi dibagi menjadi 3 module utama untuk modularitas:

1. **Auth Module**:

   - Handle autentikasi: Signup, Signin, Change Password, Refresh Token, Change Role, Signout.
   - Tabel terkait: `User`, `UserSecurity`, `ApplicationRole`, `RefreshToken`.
   - Endpoint: `/api/v1/auth/*`.
   - Integrasi: Semua endpoint lain dilindungi auth middleware.

2. **Department Module**:

   - Handle CRUD departemen (admin-only).
   - Tabel terkait: `Department`.
   - Endpoint: `/api/v1/departments/*`.
   - Tambahan: Digunakan untuk filter di attendance logs dan dashboard.

3. **Attendance Module**:
   - Handle clock in/out, list logs dengan filter dan ketepatan waktu.
   - Tabel terkait: `Attendance`, `AttendanceHistory`.
   - Endpoint: `/api/v1/attendance/*`.
   - Kompleksitas: Join dengan UserProfile dan Department untuk perhitungan punctuality (Late/On Time).

Module ini terintegrasi via shared domain models dan use cases (e.g., attendance butuh auth untuk userID).

### Persyaratan Sistem

- **Bahasa**: Go 1.22+.
- **Framework**: Fiber (web framework), GORM (ORM).
- **Database**: PostgreSQL (dengan extension `uuid-ossp` untuk UUID).
- **Dependensi Utama** (dari go.mod):
  - `github.com/gofiber/fiber/v2` (HTTP server).
  - `gorm.io/gorm` dan `gorm.io/driver/postgres` (DB ORM).
  - `github.com/go-playground/validator/v10` (validasi).
  - `github.com/sirupsen/logrus` (logging).
  - `github.com/google/uuid` (UUID generation).
  - `github.com/spf13/viper` (config .env).
  - `golang.org/x/crypto/bcrypt` (password hash).
  - Lainnya: `github.com/golang-jwt/jwt/v5` (JWT).
- **Environment Variables** (.env):
  ```
  DB_DSN=postgres://user:pass@localhost:5432/dbname?sslmode=disable
  JWT_SECRET=your-secret-key
  JWT_ACCESS_EXP=15m
  JWT_REFRESH_EXP=48h
  ```

## Cara Menjalankan Aplikasi

### 1. Instalasi Dependensi

- Clone repository: `git clone github.com/username/absensi-backend`.
- Masuk folder: `cd absensi-backend`.
- Instal Go packages: `go mod tidy`.

### 2. Setup Database

- Install PostgreSQL (local atau Docker: `docker run -p 5432:5432 -e POSTGRES_PASSWORD=pass postgres`).
- Buat database: `psql -U user -c "CREATE DATABASE dbname;"`.
- Jalankan migration: Di `main.go` atau tool terpisah, gunakan `db.AutoMigrate(&models.User{}, &models.UserProfile{}, ...)` untuk create tabel.
- Enable UUID extension: `psql -d dbname -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"`.

### 3. Konfigurasi .env

- Buat file `.env` dengan variabel di atas.

### 4. Jalankan Aplikasi

- Build: `go build -o bin/main cmd/main.go`.
- Run: `./bin/main` atau `go run cmd/main.go`.
- Server jalan di `:8080`.
- Test endpoint dengan Postman/Curl (e.g., POST `/api/v1/auth/signup` dengan body JSON).

### 5. Deployment

- Dockerize: Buat `Dockerfile` untuk build Go binary dan run.
- Deploy ke Heroku/Vercel/AWS dengan env vars.

## Endpoint API

Semua endpoint di `/api/v1`, protected by JWT kecuali auth signup/signin.

### Auth Module

- POST `/auth/signup`: Buat user (email, password, full_name) → integrasi karyawan.
- POST `/auth/signin`: Login → return access/refresh token.
- POST `/auth/change-password`: Ubah password (auth required).
- POST `/auth/refresh-token`: Refresh token.
- POST `/auth/change-role`: Ubah role (admin-only).
- POST `/auth/signout`: Logout.

### Department Module

- POST `/departments`: Create department (name, max_in, max_out) – admin-only.
- GET `/departments/:id`: Get department.
- PUT `/departments/:id`: Update department.
- DELETE `/departments/:id`: Delete department.
- GET `/departments`: List departments (pagination).

### Attendance Module

- POST `/attendance/clock-in`: Clock in (auto detect user).
- PUT `/attendance/clock-out`: Clock out.
- GET `/attendance/logs`: List logs dengan filter tanggal/departemen, ketepatan waktu, pagination.

Semua requirement soal terpenuhi: CRUD karyawan via auth/profile, CRUD departemen, absen masuk/keluar, list logs dengan perhitungan ketepatan.
