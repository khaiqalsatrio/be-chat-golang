# Rencana Implementasi Fitur Agenda & Jadwal (Agenda Plan)

Berdasarkan pembaruan di `SPEC/chat_system_spec.md`, berikut adalah rencana teknis untuk mengimplementasikan fitur pembuatan **Agenda / Event** di dalam ruang obrolan.

## 1. Konsep Dasar Fitur Agenda
Fitur ini memungkinkan anggota obrolan (baik di *Personal Chat* maupun *Group Chat*) untuk menjadwalkan acara (misalnya meeting, batas waktu/deadline, dll).
* **Fitur Utama**: Membuat, mengedit, dan menghapus Agenda.
* **Notifikasi Real-time**: Anggota grup akan mendapatkan pesan sistem saat agenda dibuat.
* **Reminder (Pengingat)**: Sistem akan mengirimkan notifikasi/pesan pengingat otomatis saat jadwal acara sudah mendekati (misalnya H-1 atau 15 menit sebelumnya).

## 2. Desain Skema Database (Prisma)
Penambahan pada skema PostgreSQL:

1. **Entitas `Agenda`**
   * `id`: UUID (Primary Key)
   * `room_id`: Relasi ke tabel `rooms`
   * `creator_id`: Relasi ke `users` (pembuat acara)
   * `title`: VARCHAR (Judul acara)
   * `description`: TEXT (Deskripsi atau link meeting)
   * `scheduled_at`: TIMESTAMP (Waktu pelaksanaan acara)
   * `created_at`: TIMESTAMP
   * `updated_at`: TIMESTAMP
2. **Entitas `AgendaParticipant` / RSVP (Opsional untuk Fase 2)**
   * `id`, `agenda_id`, `user_id`, `status` (GOING, NOT_GOING, MAYBE).

## 3. Implementasi API Endpoints (REST API)
Aksi CRUD untuk agenda dapat dilakukan via REST API standar:
* `POST /api/rooms/:roomId/agendas`: Membuat agenda baru.
* `GET /api/rooms/:roomId/agendas`: Mengambil daftar agenda di room tertentu (bisa difilter yang *upcoming* saja).
* `PUT /api/agendas/:id`: Mengedit detail agenda.
* `DELETE /api/agendas/:id`: Menghapus agenda.

## 4. Integrasi WebSocket (Real-Time Notification)
Ketika pengguna membuat atau mengedit agenda (via REST API), server backend akan langsung mengirimkan *Event Socket.io* ke room terkait.
* **Event Name**: `agenda_updated` atau `new_agenda_created`
* **Action**: Client (Frontend) yang mendengarkan event ini bisa langsung menampilkan notifikasi *"User A membuat Agenda Baru: Meeting Sinkronisasi"* secara instan.

## 5. Sistem Pengingat Otomatis (Background Jobs)
Karena kita butuh mengirim notifikasi saat waktu acara tiba (`scheduled_at`), kita perlu sistem yang berjalan di *background* untuk mengecek jadwal.
* **Teknologi**: Menggunakan *Task Queue* seperti **BullMQ** (berbasis Redis) atau **node-cron** untuk menjadwalkan pekerjaan.
* **Proses**:
  1. Saat Agenda Dibuat: Masukkan *job* baru ke *queue* dengan waktu tunggu (*delay*) sesuai jadwal acara (dikurangi waktu reminder).
  2. Saat *Job* Dieksekusi: Server secara otomatis mengirimkan pesan sistem ke dalam `rooms` via WebSocket atau FCM (*Push Notification*).

## 6. Langkah Kerja (To-Do List) Fitur Agenda
1. [ ] Update schema Prisma untuk entitas `Agenda`.
2. [ ] Jalankan migrasi database (`npx prisma migrate dev`).
3. [ ] Buat file *Controller* dan *Route* (CRUD) untuk modul Agenda.
4. [ ] Integrasikan pembuatan agenda dengan Socket.io (Emit event `new_agenda` saat agenda berhasil di-save).
5. [ ] Siapkan konfigurasi antrian/cron-job (seperti BullMQ) untuk sistem *Reminder*.
6. [ ] Tulis logika *worker* yang akan mengeksekusi pesan otomatis saat jadwal telah tiba.

---
*Fitur Agenda ini bisa dibangun secara paralel setelah fitur Group Chat dasar (dari `chat_plan.md`) selesai berfungsi.*
