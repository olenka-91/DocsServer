# 📄 DocsServer - API для хранения и управления документами

Высокопроизводительный сервер для хранения документов на Go с безопасной загрузкой файлов, JWT аутентификацией и многоуровневой системой кеширования.

##  Возможности

-  **JWT Аутентификация** - Безопасный доступ с access/refresh токенами
-  **Хранение файлов** - Эффективное хранение с кешированием на диске и в памяти
-  **Умное кеширование** - Многоуровневая стратегия кеширования для оптимальной производительности
-  **Потокобезопасность** - Конкурентный доступ с правильными механизмами блокировок
-  **Управление метаданными** - Информация о документах хранится в PostgreSQL
-  **RESTful API** - Чистый и хорошо структурированный API
-  **Docker поддержка** - Готовность к контейнеризированному развертыванию

## 🏗️ Архитектура

### Стратегия хранения
```go
// Многоуровневая система кеширования
Кеш в памяти (100MB макс) → Дисковое хранилище
```

### Особенности кеширования
- **Основано на размере**: Файлы до 2MB кешируются в памяти
- **Время жизни**: 5-минутное expiration кеша
- **LRU вытеснение**: Автоматическая очистка кеша при достижении лимитов
- **Эффективность**: Кеширование только при операциях чтения

## 🚀 Быстрый старт

### Предварительные требования
- Go 1.21+
- PostgreSQL 12+

### Установка

1. **Клонируйте репозиторий**
```bash
git clone https://github.com/olenka-91/DocsServer.git
cd DocsServer
```

2. **Настройте переменные окружения**
```bash
cp .env.example .env
# Отредактируйте .env с вашей конфигурацией
```

3. **Установите зависимости**
```bash
go mod download
```

4. **Запустите миграции базы данных**
```bash
migrate -path migrations -database "postgres://docsserver:qwerty@localhost:5436/docsserver?sslmode=disable" up 
```

5. **Запустите сервер**
```bash
go run cmd/app/main.go
```

### Развертывание с Docker

```bash
docker-compose up -d
```

## 📚 Документация API

### Эндпоинты аутентификации (публичные)

| Метод | Эндпоинт | Описание |
|-------|----------|-----------|
| `POST` | `/api/auth` | Вход пользователя (signIn) |
| `POST` | `/api/register` | Регистрация нового пользователя (signUp) |
| `POST` | `/api/refresh` | Обновление access токена (refreshToken) |

### Эндпоинты аутентификации (защищенные)

| Метод | Эндпоинт | Описание |
|-------|----------|-----------|
| `POST` | `/api/logout` | Выход пользователя (logout) |

### Эндпоинты документов (защищенные)

| Метод | Эндпоинт | Описание |
|-------|----------|-----------|
| `GET` | `/api/docs` | Получить список документов |
| `HEAD` | `/api/docs` | Получить заголовки списка документов |
| `GET` | `/api/docs/:id` | Получить документ по ID |
| `HEAD` | `/api/docs/:id` | Получить метаданные документа по ID |
| `POST` | `/api/docs` | Загрузить новый документ |
| `DELETE` | `/api/docs/:id` | Удалить документ по ID |

### Примеры использования

**Регистрация пользователя:**
```bash
curl -L -X POST 'localhost:8000/api/register/' \
-H 'Content-Type: application/json' \
-d '{
    "name": "OlgaDvornikova7",
    "password": "Uprising123_"
}'
```

**Вход пользователя:**
```bash
curl -L -X POST 'localhost:8000/api/auth/' \
-H 'Content-Type: application/json' \
-d '{
    "name": "OlgaDvornikova7",
    "password": "Uprising123_"
}'
```

**Загрузка документа:**
```bash
curl -L -X POST 'localhost:8000/api/docs' \
 -H "Authorization: Bearer <ваш_jwt_токен>" \
 -F 'meta=" {
      \"name\": \"file.jpg\",
      \"file\": true,
      \"public\": true,      
      \"mime\": \"image/jpg\",
      \"grant\": [\"FirstUser3\", \"OlgaDvornikova7\"]
    }"' \
-F 'file=@"/C:/Users/Ольга/Desktop/og_og.jpg"'
```

**Получение списка документов:**
```bash
curl -L -X GET 'localhost:8000/api/docs/?key=filename&value=photo&limit=10' \
  -H "Authorization: Bearer <ваш_jwt_токен>"
```

**Скачивание документа:**
```bash
curl -L -X GET 'localhost:8000/api/docs/7ea0a0b8-c652-41a4-86de-678a0e214c8c' \
-H 'Authorization:  Bearer <ваш_jwt_токен>'
```

**Удаление документа:**
```bash
curl -L -X DELETE 'localhost:8000/api/docs/7ea0a0b8-c652-41a4-86de-678a0e214c8c' \
  -H "Authorization: Bearer <ваш_jwt_токен>"
```

## ⚙️ Конфигурация

### Переменные окружения

```env
# Сервер
SERVER_PORT=8080
SERVER_HOST=localhost

# База данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=docsserver

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Хранилище
STORAGE_PATH=./storage
MAX_FILE_SIZE=10485760  # 10MB
```

### Конфигурация кеша

```go
const (
    cacheTTL             = 5 * time.Minute
    maxMemoryCacheSize   = 100 * 1024 * 1024  // 100MB
    maxMemoryCachedFiles = 100
    maxFileCacheSize     = 2 * 1024 * 1024    // 2MB на файл
)
```

## 🏛️ Структура проекта

```
DocsServer/
├── cmd/
│   ├── app/          # Основное приложение
│  
├── internal/
│   ├── config/          # Управление конфигурацией
│   ├── handler/         # HTTP контроллеры
      └── middleware/    # HTTP middleware
│   ├── entity/          # Сущности базы данных
│   ├── utils/           # Функции для работы с токеном и паролем
│   ├── repository/      # Уровень доступа к данным
│   ├── service/         # Бизнес-логика
│   └── storage/         # Реализация файлового хранилища
├── pkg/
│   ├── httpserver/      # Реализация сервера
└── migrate/             # Миграции базы данных
```

## 🛠️ Используемые технологии

- **Фреймворк**: [Gin](https://github.com/gin-gonic/gin)
- **База данных**: [PostgreSQL](https://www.postgresql.org/) 
- **Аутентификация**: [JWT](https://github.com/golang-jwt/jwt)
- **Логирование**: [Logrus](https://github.com/sirupsen/logrus)
- **Конфигурация**: [Viper](https://github.com/spf13/viper) (опционально)
  
## 👥 Авторы

- **Olenka-91** - Начальная работа - [GitHub](https://github.com/olenka-91)

---

**⭐ Поставьте звезду репозиторию, если он вам полезен!**
