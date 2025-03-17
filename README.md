# Lower Underscore Case

Lower Underscore Case adalah aplikasi backend yang dibangun menggunakan **Golang** dengan framework **Gin**. Aplikasi ini menyediakan fitur autentikasi menggunakan JWT serta operasi CRUD untuk entitas tertentu.

## Fitur

- **Autentikasi**: Login, Logout, dan Verifikasi Token JWT.
- **CRUD**: Operasi Create, Read, Update, dan Delete untuk entitas utama.
- **Middleware**: Middleware untuk manajemen autentikasi dan otorisasi.
- **Dokumentasi API**: Dokumentasi API menggunakan Swagger.

## Instalasi & Menjalankan Proyek

### 1. Clone Repository

```bash
git clone https://github.com/charlospurba/lower_underscore_case.git
cd lower_underscore_case
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Konfigurasi Environment

Buat file `.env` di root direktori proyek dan tambahkan konfigurasi berikut:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
JWT_SECRET=your_jwt_secret
```

### 4. Migrasi Database

Pastikan database sudah dibuat dan dapat diakses. Kemudian, jalankan migrasi:

```bash
go run migrations/migrate.go
```

### 5. Generate Dokumentasi Swagger

Instal Swagger terlebih dahulu jika belum:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Kemudian, generate dokumentasi:

```bash
swag init
```

### 6. Jalankan Server

```bash
go run main.go
```

Server akan berjalan di: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

Note:
ERD -> https://dbdiagram.io/d/67cd9fa6263d6cf9a0bd588a

