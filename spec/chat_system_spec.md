# Spesifikasi Sistem Backend Chat Real-Time

Dokumen ini mendefinisikan spesifikasi dan arsitektur untuk sistem backend chat real-time yang kompleks.

## 1. Fitur Utama (Core Features)
* **Real-Time Communication**: Menggunakan protokol WebSockets (misalnya Socket.io atau ws) untuk komunikasi dua arah dengan latensi rendah.
* **1-on-1 Chat (Direct Messaging)**: Percakapan pribadi antar dua pengguna.
* **Group Chat**: Percakapan grup dengan banyak anggota, mencakup manajemen grup (Admin, Member, Kick, Add, dll).
* **Message Status / Read Receipts**: Melacak status pesan (Sent, Delivered, Read).
* **User Presence**: Melacak status pengguna secara real-time (Online, Offline, Last Seen).
* **Typing Indicators**: Menampilkan status saat pengguna sedang mengetik.

## 2. Fitur Lanjutan (Advanced Features)
* **Multimedia Attachments**: Dukungan untuk mengirim file, gambar, video, dokumen, dan voice notes (pesan suara).
* **Pesan Bersarang (Replies) & Forwarding**: Membalas pesan spesifik atau meneruskan pesan ke obrolan lain.
* **Reactions**: Memberikan reaksi emoji pada pesan spesifik.
* **Edit & Hapus Pesan**: Pengguna dapat mengedit atau menghapus pesan (Soft delete / Hard delete) dengan batas waktu tertentu.
* **Pinned Messages**: Menyematkan pesan penting di dalam grup atau personal chat.
* **Push Notifications**: Notifikasi offline/background menggunakan layanan seperti Firebase Cloud Messaging (FCM).
* **Pencarian Pesan**: Pencarian teks penuh (Full-text search) pada riwayat obrolan.
* **Integrasi AI (Artificial Intelligence)**: Fitur asisten cerdas yang terintegrasi ke dalam chat (mirip Meta AI) untuk berinteraksi langsung, merangkum percakapan, atau memberikan saran balasan (*smart replies*).
* **Fitur Agenda / Jadwal**: Memungkinkan pembuatan jadwal pertemuan, kalender event, atau *reminder* di dalam grup maupun percakapan personal, lengkap dengan sistem notifikasi otomatis.

## 3. Arsitektur Teknis
* **Protokol Utama**: 
  - REST API (untuk autentikasi, manajemen profil, riwayat pesan, dll).
  - WebSockets (untuk real-time events seperti kirim pesan, status baca, typing).
* **Database Relasional / NoSQL**:
  - PostgreSQL / MySQL (untuk data user, relasi grup, metadata).
  - MongoDB (opsional, cocok untuk menyimpan log pesan dalam jumlah besar yang fleksibel).
* **Caching & Pub/Sub**:
  - Redis: Digunakan untuk session management, menyimpan status Online/Offline pengguna secara cepat, dan sebagai Message Broker (Pub/Sub) jika sistem di-scale ke beberapa server websocket.
* **File Storage**:
  - AWS S3, Google Cloud Storage, atau MinIO untuk menyimpan file attachment (gambar, dokumen, dll).
* **Layanan AI (AI Services)**:
  - Menggunakan API LLM pihak ketiga (seperti OpenAI GPT, Google Gemini, atau Anthropic Claude) untuk memproses teks/perintah dari user dan menghasilkan respons secara *real-time*.

## 4. Model Data Dasar (Entitas)
### a. User
* `id`, `username`, `email`, `password_hash`, `avatar_url`, `status`, `last_seen`, `created_at`
### b. Chat Room / Conversation
* `id`, `type` (DIRECT, GROUP), `name` (untuk grup), `created_at`, `updated_at`
### c. Chat Participant
* `room_id`, `user_id`, `role` (ADMIN, MEMBER), `joined_at`
### d. Message
* `id`, `room_id`, `sender_id` (Bisa ID User atau ID khusus AI Agent), `message_type` (TEXT, IMAGE, VIDEO, FILE, VOICE), `content` (teks atau URL), `reply_to_id`, `is_ai_generated` (boolean), `created_at`, `updated_at`, `deleted_at`
### e. Message Status / Receipt
* `message_id`, `user_id`, `status` (DELIVERED, READ), `timestamp`
### f. Agenda / Event
* `id`, `room_id`, `creator_id`, `title`, `description`, `scheduled_at` (Waktu acara), `created_at`, `updated_at`

## 5. Event WebSockets Utama
### a. Client to Server (Emit)
* `join_room` / `leave_room`: Masuk/keluar dari scope obrolan tertentu.
* `send_message`: Mengirim pesan baru.
* `typing_start` / `typing_end`: Trigger saat mulai/selesai mengetik.
* `mark_as_read`: Menandai pesan telah dibaca.
* `edit_message` / `delete_message`: Modifikasi pesan.

### b. Server to Client (Listen)
* `receive_message`: Menerima pesan baru secara real-time.
* `message_status_update`: Update visual ketika pesan berubah jadi Delivered/Read.
* `user_typing`: Menampilkan "User is typing...".
* `presence_update`: Update status online teman/kontak.

## 6. Keamanan & Skalabilitas
* **Authentication**: JWT (JSON Web Tokens) untuk REST API dan koneksi WebSocket awal.
* **Rate Limiting**: Mencegah spam pesan dalam waktu singkat.
* **Pagination**: Memuat riwayat obrolan secara bertahap (Cursor-based pagination disarankan untuk chat).
* **Scaling**: Memisahkan server WebSocket (Stateful) dan REST API server (Stateless). Menggunakan Redis Pub/Sub agar pesan bisa dikirim lintas instance WebSocket.

---
*Spesifikasi ini dapat disesuaikan kembali sesuai kebutuhan spesifik project saat fase implementasi dimulai.*
