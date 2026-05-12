# Backend Reconstruction Plan: Chat Conversations & Last Message

Dokumen ini berisi panduan teknis untuk memperbaiki endpoint Backend agar sinkron dengan kebutuhan tampilan aplikasi mobile, khususnya untuk memunculkan pesan terakhir di daftar percakapan.

## 1. Perbaikan Endpoint `GET /api/chat/conversations`

Saat ini, aplikasi mobile mengharapkan setiap objek percakapan memiliki ringkasan pesan terakhir.

### Struktur Response yang Diharapkan (JSON)
Pastikan field `last_message` dikirimkan sebagai string langsung di dalam objek room/conversation.

```json
{
  "status": "success",
  "data": [
    {
      "id": "room-uuid",
      "name": "Username",
      "type": "PRIVATE",
      "last_message": "Isi pesan terakhir di sini...",
      "updated_at": "2023-12-01T15:00:00Z",
      "participants": [
        {
          "id": "user-uuid",
          "username": "dona",
          "avatar_url": "https://..."
        }
      ]
    }
  ]
}
```

## 2. Implementasi di Sisi Go (Backend)

### A. Definisi Struct
Pastikan struct yang digunakan untuk Response memiliki tag JSON `last_message`.

```go
type ConversationResponse struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`
    LastMessage string    `json:"last_message"` // Tambahkan tag ini
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### B. Logic Query (Database)
Jika menggunakan SQL/GORM, gunakan subquery atau join untuk mengambil pesan terbaru secara efisien.

**Contoh Logic SQL:**
```sql
SELECT 
    rooms.*, 
    (SELECT content FROM messages 
     WHERE messages.room_id = rooms.id 
     ORDER BY created_at DESC LIMIT 1) as last_message
FROM rooms
...
```

## 3. Sinkronisasi Real-time

### A. Update `updated_at` Room
Sangat penting untuk mengupdate field `updated_at` pada tabel `rooms` setiap kali ada pesan baru (`SendMessage`). Ini memastikan urutan chat di mobile selalu yang terbaru berada di atas.

```go
// Setiap kali pesan dikirim:
updateRoomTime := "UPDATE rooms SET updated_at = NOW() WHERE id = ?"
```

### B. WebSocket Payload Consistency
Pastikan event WebSocket `new_message` mengirimkan field yang konsisten dengan HTTP API:
- Gunakan `room_id` atau `conversation_id`.
- Gunakan field `content` untuk teks pesan.

## 4. Checklist Verifikasi
- [ ] Endpoint `/api/chat/conversations` mengembalikan field `last_message`.
- [ ] Field `last_message` bukan berupa objek kosong/null jika sudah ada chat.
- [ ] Urutan list conversation berdasarkan `updated_at` terbaru (DESC).
- [ ] Nama field menggunakan snake_case (`last_message`) agar sesuai dengan mapping Frontend saat ini.

---
**Catatan:** Frontend saat ini sudah dikonfigurasi untuk membaca field `last_message`. Begitu Backend mengirimkan data ini, teks "No messages yet" akan otomatis terganti dengan isi pesan asli.
