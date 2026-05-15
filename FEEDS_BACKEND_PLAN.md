# Perencanaan Backend - Fitur Feeds Lengkap

Dokumen ini merincikan rencana pengembangan fitur Feeds pada backend menggunakan Go (Clean Architecture).

## 1. Skema Database (PostgreSQL)

### Table: `posts`
Menyimpan data postingan utama.
- `id`: UUID (Primary Key)
- `user_id`: UUID (Foreign Key ke `users.id`)
- `caption`: TEXT
- `image_url`: VARCHAR(255)
- `created_at`: TIMESTAMP
- `updated_at`: TIMESTAMP

### Table: `likes`
Menyimpan interaksi "Suka" pada postingan.
- `id`: UUID (Primary Key)
- `post_id`: UUID (Foreign Key ke `posts.id`)
- `user_id`: UUID (Foreign Key ke `users.id`)
- `created_at`: TIMESTAMP

### Table: `comments`
Menyimpan komentar pada postingan.
- `id`: UUID (Primary Key)
- `post_id`: UUID (Foreign Key ke `posts.id`)
- `user_id`: UUID (Foreign Key ke `users.id`)
- `content`: TEXT
- `created_at`: TIMESTAMP

---

## 2. Struktur Folder & Entitas (Domain Layer)

### [NEW] `src/internal/domain/entities/post.go`
```go
type Post struct {
    ID        uuid.UUID `json:"id"`
    UserID    uuid.UUID `json:"user_id"`
    User      User      `json:"user"` // Relasi ke entitas User
    Caption   string    `json:"caption"`
    ImageURL  string    `json:"image_url"`
    LikesCount int      `json:"likes_count"`
    CommentsCount int   `json:"comments_count"`
    IsLiked    bool     `json:"is_liked"` // Berdasarkan konteks user saat ini
    CreatedAt time.Time `json:"created_at"`
}
```

### [NEW] `src/internal/domain/entities/comment.go`
```go
type Comment struct {
    ID        uuid.UUID `json:"id"`
    PostID    uuid.UUID `json:"post_id"`
    UserID    uuid.UUID `json:"user_id"`
    User      User      `json:"user"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}
```

---

## 3. Repositori (Data Layer)

### [NEW] `src/internal/domain/repositories/post_repository.go`
Definisi interface untuk interaksi database:
- `Create(post *entities.Post) error`
- `FindAll(limit, offset int, currentUserID uuid.UUID) ([]*entities.Post, error)`
- `FindByID(id uuid.UUID) (*entities.Post, error)`
- `Delete(id uuid.UUID) error`
- `ToggleLike(postID, userID uuid.UUID) (bool, error)`
- `AddComment(comment *entities.Comment) error`

---

## 4. Use Cases (Business Logic Layer)

### [NEW] `src/internal/usecases/post_usecase.go`
Logika bisnis untuk:
1. **CreatePost**: Memproses input postingan baru dan menyimpan ke DB.
2. **GetGlobalFeed**: Mengambil daftar postingan terbaru dengan pagination.
3. **ToggleLike**: Menambah atau menghapus "like" (logika: jika sudah ada, hapus; jika belum, tambah).
4. **PostComment**: Validasi dan simpan komentar baru.

---

## 5. Interface / API Endpoints (HTTP Layer)

### [NEW] `src/internal/interfaces/http/handlers/post_handler.go`

Endpoints:
- `POST /api/posts` (Butuh Auth)
    - Input: `caption`, `image` (multipart form atau URL)
- `GET /api/posts` (Butuh Auth)
    - Query: `limit`, `offset`
- `POST /api/posts/:id/like` (Butuh Auth)
    - Toggle status suka.
- `POST /api/posts/:id/comments` (Butuh Auth)
    - Input: `content`
- `GET /api/posts/:id/comments`
    - Daftar komentar untuk sebuah postingan.
- `DELETE /api/posts/:id` (Butuh Auth)
    - Hapus postingan sendiri.

---

## 6. Rencana Verifikasi
1. **Unit Testing**: Mengetes Repository dan Usecase menggunakan mock database.
2. **Integration Testing**: Mengetes API menggunakan Postman atau alat sejenis.
3. **Mobile Integration**: Menghubungkan UI React Native yang sudah dibuat dengan API ini.
