# Rencana Implementasi Autentikasi (Auth Plan)

Berdasarkan spesifikasi di `SPEC/chat_system_spec.md`, berikut adalah rencana bertahap untuk membangun sistem autentikasi.

## 1. Persiapan Lingkungan & Teknologi (Tech Stack)
Berdasarkan kesepakatan, kita akan menggunakan *stack* teknologi berikut:
* **Bahasa**: Go (Golang)
* **Framework Backend**: NestJS
* **Database**: PostgreSQL
* **ORM**: TypeORM
* **Security & Auth Tools**: 
  - `bcrypt` untuk hashing password.
  - `@nestjs/jwt` untuk membuat dan memvalidasi token.
  - `passport`, `passport-jwt`, `passport-google-oauth20` untuk strategi autentikasi.
  - `zod` atau `class-validator` untuk validasi input data.

## 2. Pembuatan Model / Skema Database
Berdasarkan spesifikasi dasar, kita perlu mendefinisikan skema tabel `users`:
```sql
-- Contoh struktur tabel (representasi logis)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255), -- Opsional (null untuk user yang mendaftar via Google)
    google_id VARCHAR(255) UNIQUE, -- Menyimpan ID dari Google (Opsional)
    avatar_url VARCHAR(255),
    status VARCHAR(20) DEFAULT 'OFFLINE', -- ONLINE, OFFLINE
    last_seen TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 3. Implementasi API Endpoints

### a. Endpoint Registrasi (`POST /api/auth/register`)
* **Input**: `username`, `email`, `password`.
* **Proses**:
  1. Validasi format input (email valid, panjang password, dll).
  2. Cek apakah `email` atau `username` sudah terdaftar.
  3. Hash `password` menggunakan bcrypt (misal dengan 10 salt rounds).
  4. Simpan data user baru ke database.
* **Output**: Status 201 Created dengan pesan sukses.

### b. Endpoint Login (`POST /api/auth/login`)
* **Input**: `email`, `password`.
* **Proses**:
  1. Cari user berdasarkan `email`. Jika tidak ada, kembalikan error "Kredensial tidak valid".
  2. Verifikasi/cocokkan `password` input dengan `password_hash` dari database menggunakan bcrypt.
  3. Jika valid, buat dan tandatangani JWT (Access Token).
     * *Payload JWT disarankan berisi*: `userId`, `username`.
     * *Masa berlaku (Expiration)*: misal 1 jam (1h) atau 1 hari (1d).
* **Output**: Status 200 OK, mengembalikan `{ "token": "eyJh...", "user": { "id": "...", "username": "..." } }`.

### c. Endpoint Profil Pengguna (`GET /api/auth/me` atau `/api/users/me`)
* **Tujuan**: Untuk mengambil data profil user yang sedang login, dan memverifikasi apakah token klien masih valid.
* **Proses**:
  1. Di-protect oleh *Auth Middleware*.
  2. Mengambil data user (kecuali password) berdasarkan `userId` dari token.
* **Output**: Detail pengguna saat ini.

### d. Endpoint Google Login (`POST /api/auth/google`)
* **Input**: `idToken` (Token Google yang didapat dari sisi Frontend/Aplikasi Client).
* **Proses**:
  1. Verifikasi `idToken` ke server Google menggunakan `google-auth-library`.
  2. Ekstrak informasi profil (`email`, `name`, `picture`) dari token yang valid.
  3. Cek apakah user dengan `email` tersebut sudah terdaftar di database.
     * Jika belum: Buat akun baru secara otomatis (kosongkan `password_hash`, simpan `google_id` dan `avatar_url`).
     * Jika sudah: Gunakan akun yang ada.
  4. Generate JWT (Access Token) untuk user tersebut sebagai tanda login berhasil di sistem kita.
* **Output**: Status 200 OK, mengembalikan `{ "token": "eyJh...", "user": {...} }`.

## 4. Pembuatan Auth Middleware
Sebuah *middleware* (misal `verifyToken`) yang akan digunakan untuk melindungi endpoint REST API yang membutuhkan login, serta nantinya untuk melindungi koneksi WebSocket.
* **Proses**:
  1. Mengekstrak token dari Header `Authorization: Bearer <token>`.
  2. Memverifikasi keaslian token menggunakan `JWT_SECRET`.
  3. Menyisipkan data `userId` ke dalam *request object* agar bisa diakses oleh *controller* berikutnya.

## 5. Langkah Kerja (To-Do List) Berikutnya
Langkah eksekusi yang akan kita lakukan selanjutnya adalah:
1. [ ] Inisialisasi project backend (`npm init`, setup TypeScript).
2. [ ] Install dependencies (`express`, `prisma`, `pg`, `jsonwebtoken`, `bcrypt`, dll).
3. [ ] Inisialisasi Prisma dan setup skema Database PostgreSQL untuk entitas `User`.
4. [ ] Buat file *Controller* dan *Route* untuk Register, Login Tradisional, & **Google Login**.
5. [ ] Buat dan terapkan *Auth Middleware*.
6. [ ] Uji API menggunakan Postman, cURL, atau REST Client.

---
*Jika Anda sudah siap, beri tahu saya untuk mulai mengeksekusi langkah pertama (Inisialisasi project)!*
