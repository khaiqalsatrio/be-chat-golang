# Backend Implementation Plan: Status/Story Feature

Dokumen ini merinci langkah-langkah yang diperlukan untuk memperbarui Backend (Go) guna mendukung fitur Status/Story.

## 1. Skema Database (PostgreSQL/MySQL)

Kita membutuhkan tabel baru bernama `statuses` untuk menyimpan data kiriman pengguna.

```sql
CREATE TABLE statuses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    media_url TEXT NOT NULL,
    media_type VARCHAR(20) DEFAULT 'image', -- 'image' atau 'video'
    caption TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL -- created_at + 24 hours
);

-- Index untuk performa fetch status yang aktif
CREATE INDEX idx_active_status ON statuses (expires_at) WHERE expires_at > CURRENT_TIMESTAMP;
```

## 2. API Endpoints

Berikut adalah daftar endpoint yang perlu diimplementasikan:

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/status` | Mengunggah status baru (mendukung upload file) |
| `GET` | `/api/status` | Mengambil daftar status terbaru dari kontak/teman |
| `GET` | `/api/status/me` | Mengambil riwayat status milik sendiri |
| `DELETE` | `/api/status/:id` | Menghapus status sebelum masa kadaluwarsa |

## 3. Alur Kerja (Workflow)

### A. Pengunggahan (Upload)
1. Frontend mengirim file gambar/video melalui `multipart/form-data`.
2. Backend memproses file dan menyimpannya di Cloud Storage (S3/Cloudinary/Local Storage).
3. Simpan metadata (URL, user_id, expires_at) ke database.

### B. Pengambilan Data (Fetching)
1. Backend melakukan query untuk mencari status yang:
   - Milik pengguna yang berteman dengan kita.
   - `expires_at` masih di masa depan (> waktu sekarang).
2. Mengelompokkan status berdasarkan `user_id` agar frontend bisa menampilkannya dalam bentuk Story Circles.

## 4. Sistem Kadaluwarsa (Automatic Expiry)

Meskipun kita hanya menampilkan status yang belum kadaluwarsa lewat filter query, kita tetap butuh "pembersihan" berkala untuk menghemat storage:

- **Cron Job / Worker**: Jalankan tugas setiap jam untuk menghapus data atau menandai status yang sudah lewat dari 24 jam sebagai `archived`.

## 5. Integrasi Frontend (Next Step)

Setelah API siap, frontend akan:
- Mengganti `DUMMY_STATUS` di `StatusList.tsx` dengan data asli dari API.
- Menambahkan fitur kamera/galeri untuk memicu `POST /api/status`.

---
*Dibuat oleh: Antigravity AI Coding Assistant*
