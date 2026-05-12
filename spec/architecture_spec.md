src/
├── cmd/
│   └── server/
│       └── main.go                # Entry point aplikasi
│
├── config/                        # Konfigurasi ENV & Constants
│   ├── config.go
│   └── env.go
│
├── internal/
│   ├── domain/                    # Layer Domain
│   │   ├── entities/
│   │   │   ├── user.go
│   │   │   ├── message.go
│   │   │   └── room.go
│   │   │
│   │   ├── repositories/
│   │   │   ├── user_repository.go
│   │   │   ├── chat_repository.go
│   │   │   └── room_repository.go
│   │   │
│   │   └── services/
│   │       ├── ai_service.go
│   │       └── auth_service.go
│   │
│   ├── usecases/                  # Business Logic
│   │   ├── auth/
│   │   │   ├── login_usecase.go
│   │   │   └── register_usecase.go
│   │   │
│   │   ├── chat/
│   │   │   ├── send_message.go
│   │   │   ├── get_messages.go
│   │   │   └── ai_chat.go
│   │   │
│   │   └── room/
│   │       └── create_room.go
│   │
│   ├── interfaces/                # Interface Adapters
│   │   ├── http/
│   │   │   ├── handlers/
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── chat_handler.go
│   │   │   │   └── room_handler.go
│   │   │   │
│   │   │   ├── middleware/
│   │   │   │   ├── auth_middleware.go
│   │   │   │   └── logger_middleware.go
│   │   │   │
│   │   │   ├── dto/
│   │   │   │   ├── auth_request.go
│   │   │   │   └── chat_request.go
│   │   │   │
│   │   │   └── routes/
│   │   │       └── routes.go
│   │   │
│   │   └── websocket/
│   │       ├── ws_handler.go
│   │       └── ws_hub.go
│   │
│   └── infrastructure/            # Layer Infrastruktur
│       ├── database/
│       │   ├── postgres.go
│       │   └── migrations/
│       │
│       ├── repositories/
│       │   ├── postgres_user_repository.go
│       │   ├── postgres_chat_repository.go
│       │   └── postgres_room_repository.go
│       │
│       ├── services/
│       │   ├── jwt_service.go
│       │   ├── openai_service.go
│       │   ├── google_auth_service.go
│       │   └── redis_service.go
│       │
│       └── logger/
│           └── zap_logger.go
│
├── pkg/                           # Shared reusable packages
│   ├── response/
│   │   └── response.go
│   │
│   ├── validator/
│   │   └── validator.go
│   │
│   └── utils/
│       └── helper.go
│
├── go.mod
├── go.sum
└── .env