# Chat App API Documentation

This document provides a comprehensive guide to the Chat Application API, including authentication, available endpoints, and the WebSocket protocol.

## Base URL
`http://localhost:8080/api`

---

## Authentication
Most endpoints require a JWT Bearer token. Include the token in the `Authorization` header:
```http
Authorization: Bearer <your_jwt_token>
```

---

## 1. Authentication Endpoints (`/auth`)

### **Register User**
*   **URL:** `/auth/register`
*   **Method:** `POST`
*   **Success Response (201):** Returns user profile and JWT token.

### **Login User**
*   **URL:** `/auth/login`
*   **Method:** `POST`
*   **Success Response (200):** Returns JWT token.

---

## 2. Chat & Conversations (`/chat`) - Recommended for Mobile

### **Get Conversations**
List all rooms the current user is a participant of.
*   **URL:** `/chat/conversations`
*   **Method:** `GET`

### **Get Conversation Messages**
Fetch message history for a specific conversation.
*   **URL:** `/chat/conversations/:roomId/messages`
*   **Method:** `GET`
*   **Query Params:** `limit`, `offset`

### **Send Message**
Send a message via HTTP (will trigger real-time broadcast).
*   **URL:** `/chat/messages`
*   **Method:** `POST`
*   **Body:**
    ```json
    {
      "conversation_id": "room-uuid",
      "content": "Hello world!",
      "type": "TEXT"
    }
    ```

### **Create Conversation**
*   **URL:** `/chat/conversations`
*   **Method:** `POST`
*   **Body:**
    ```json
    {
      "type": "PRIVATE",
      "participants": ["user-uuid"]
    }
    ```

---

## 3. Legacy Room Endpoints (`/rooms`)
*   `GET /rooms`
*   `POST /rooms`
*   `GET /rooms/:roomId/messages`
*   `POST /rooms/:roomId/messages`

---

## 4. WebSocket Protocol (`/ws`)

Real-time updates for new messages.

### **Connection**
`ws://localhost:8080/api/ws?token=<your_jwt_token>`

### **Real-time Event**
When a new message is sent (via HTTP or WS), all connected clients in the room will receive a JSON payload:
```json
{
  "id": "message-uuid",
  "room_id": "room-uuid",
  "sender_id": "sender-uuid",
  "content": "Hello world!",
  "type": "TEXT",
  "sender": {
    "id": "sender-uuid",
    "username": "johndoe",
    "avatar_url": "..."
  },
  "created_at": "2023-12-01T10:00:00Z"
}
```

---

## Entities Reference

### **User**
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "avatar_url": "string",
  "status": "ONLINE/OFFLINE"
}
```

### **Message**
```json
{
  "id": "uuid",
  "room_id": "uuid",
  "sender_id": "uuid",
  "content": "string",
  "type": "TEXT",
  "sender": { ...User... },
  "created_at": "datetime"
}
```
