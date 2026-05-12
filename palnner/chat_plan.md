# Rencana Implementasi Fitur Chat Real-Time (Chat Plan)

Berdasarkan spesifikasi di `SPEC/chat_system_spec.md`, dokumen ini merinci rencana implementasi untuk fitur inti yaitu **Real-Time Chat** menggunakan WebSockets.

## 1. Persiapan Teknologi Real-Time
Untuk komunikasi real-time yang stabil di atas Express.js, kita akan menggunakan **Socket.io**.
*   **Keunggulan Socket.io**: Memiliki fitur built-in untuk *Rooms* (sangat berguna untuk group chat / private chat), otomatis *reconnect*, dan mudah di-scale menggunakan Redis Adapter jika ke depannya ada lebih dari satu server.

## 2. Desain Skema Database (Prisma / PostgreSQL)
Kita perlu menambahkan tabel-tabel baru ke dalam skema database yang sudah ada:

1.  **Room / Conversation (`rooms`)**
    *   `id` (UUID), `name` (Opsional, untuk grup), `type` (DIRECT, GROUP), `created_at`.
2.  **Room Participants (`room_participants`)**
    *   `id`, `room_id`, `user_id`, `role` (ADMIN/MEMBER), `joined_at`.
3.  **Messages (`messages`)**
    *   `id`, `room_id`, `sender_id`, `content`, `message_type` (TEXT, IMAGE, dll), `created_at`, `updated_at`.
4.  **Message Status (`message_status`)** -> *Untuk fitur Read Receipts*
    *   `id`, `message_id`, `user_id`, `status` (DELIVERED, READ), `updated_at`.

## 3. Alur Autentikasi WebSocket
Koneksi WebSocket juga harus diotentikasi agar tidak sembarang orang bisa terhubung.
*   **Proses**: Saat *client* mencoba terhubung ke Socket.io, mereka harus mengirimkan JWT Token (bisa lewat *handshake auth*).
*   **Middleware Socket.io**: Akan memverifikasi token. Jika valid, koneksi diterima dan `socket.user` akan menyimpan data pengguna. Jika tidak valid, koneksi ditolak.

## 4. Alur Event WebSocket (Client & Server)

### A. Manajemen Room (Masuk & Keluar)
*   **Client Emit**: `join_room` membawa data `roomId`.
*   **Server Action**: Memverifikasi apakah user merupakan partisipan dari `roomId` tersebut di database. Jika ya, masukkan socket ke dalam *room* Socket.io menggunakan `socket.join(roomId)`.

### B. Mengirim & Menerima Pesan
*   **Client Emit**: `send_message` membawa data `{ roomId, content, type }`.
*   **Server Action**:
    1.  Validasi data.
    2.  Simpan pesan ke database PostgreSQL (Tabel `messages`).
    3.  *Broadcast* pesan ke semua user yang ada di *room* tersebut menggunakan `io.to(roomId).emit('receive_message', new_message)`.

### C. Indikator Mengetik (Typing)
*   **Client Emit**: `typing_start` / `typing_end` membawa data `roomId`.
*   **Server Action**: Mem-forward event ini ke user lain di room yang sama menggunakan `socket.to(roomId).emit('user_typing', { userId })`. (Tidak perlu disimpan di database).

### D. Status Pesan (Read Receipts)
*   **Client Emit**: `mark_as_read` membawa data `messageId` dan `roomId`.
*   **Server Action**:
    1.  Update tabel `message_status` di database menjadi `READ`.
    2.  *Broadcast* ke *room* (terutama ke pengirim) dengan `io.to(roomId).emit('message_status_update', { messageId, status: 'READ' })`.

## 5. REST API Pendukung
Selain WebSocket, kita tetap membutuhkan REST API untuk aksi-aksi *stateless*:
*   `GET /api/chats`: Mengambil daftar obrolan (Room) yang diikuti user.
*   `GET /api/chats/:roomId/messages`: Mengambil riwayat pesan dari obrolan tertentu (*Pagination/Cursor-based*).
*   `POST /api/chats/group`: Membuat obrolan grup baru.

## 6. Langkah Kerja (To-Do List) Fitur Chat
1. [ ] Install dependency real-time (`socket.io`).
2. [ ] Update Prisma Schema untuk entitas `Room`, `Participant`, `Message`, dan `MessageStatus`.
3. [ ] Jalankan migrasi database (`npx prisma migrate dev`).
4. [ ] Inisialisasi Server Socket.io yang terintegrasi dengan server HTTP Express.
5. [ ] Buat Middleware Autentikasi khusus untuk Socket.io.
6. [ ] Implementasi event handlers: `join_room`, `send_message`, dan `typing`.
7. [ ] Buat REST API *Controllers* untuk mengambil daftar chat dan riwayat pesan.
8. [ ] Uji coba komunikasi dua arah (bisa menggunakan Postman WebSocket client atau script frontend sederhana).

---
*Dokumen ini merupakan panduan spesifik untuk fitur chat. Fitur ini akan kita eksekusi setelah fondasi proyek dan Autentikasi (Auth) selesai dibangun.*
