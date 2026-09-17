# Hotel Ops

Учебное full-stack приложение для управления гостиницей: инциденты, сотрудники, номера и бронирования. Есть роли пользователей, JWT-аутентификация и email-подтверждения бронирований.

## Стек

Go (Chi, pgx, bcrypt, JWT), PostgreSQL, React/Vite/Axios, Docker Compose и GitHub Actions.

## Быстрый запуск

Требуется Docker Desktop.

```bash
docker compose up --build
```

Перед первым запуском создайте `.env` из `.env.example` (`Copy-Item .env.example .env` в PowerShell или `cp .env.example .env` в macOS/Linux).

- SPA для сотрудников: `http://localhost:5173`
- Публичное бронирование: `http://localhost:5174`
- API: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/swagger/`
- OpenAPI: `http://localhost:8080/openapi.yaml`

PostgreSQL, миграции и демонстрационные данные запускаются автоматически. `.env` не коммитится; перед публикацией замените `JWT_SECRET` на случайное значение.

## Конфигурация

| Переменная | Назначение |
| --- | --- |
| `JWT_SECRET` | Секрет подписи JWT. |
| `FRONTEND_ORIGIN` | Разрешённый origin для CORS. |
| `SMTP_HOST`, `SMTP_PORT` | SMTP-сервер и порт. |
| `SMTP_USERNAME`, `SMTP_PASSWORD` | Учётные данные SMTP / пароль приложения. |
| `SMTP_FROM` | Адрес отправителя писем. |

Если SMTP не настроен, бронирования создаются, но письмо гостю не отправляется.

## API

Все ответы API — JSON. Ошибки: `{"error":{"code":"...","message":"..."}}`. Защищённые маршруты требуют `Authorization: Bearer <access_token>`; роли проверяются сервером. Полный интерактивный контракт — в Swagger.

| Метод | Маршрут | Назначение |
| --- | --- | --- |
| POST | `/api/v1/auth/register`, `/login`, `/refresh`, `/logout` | Регистрация и управление сессией. |
| GET | `/api/v1/auth/me` | Текущий пользователь. |
| GET | `/api/v1/public/hotels`, `/public/rooms` | Публичный каталог. |
| POST | `/api/v1/public/bookings` | Создание бронирования. |
| GET, POST | `/api/v1/incidents` | Инциденты. |
| GET, PUT, DELETE | `/api/v1/incidents/{id}` | Один инцидент. |
| GET, POST | `/api/v1/employees` | Сотрудники. |
| GET, PUT, DELETE | `/api/v1/employees/{id}` | Один сотрудник. |
| GET, POST, DELETE | `/api/v1/hotels`, `/categories`, `/locations` | Справочники. |
| GET, POST | `/api/v1/admin/room-types`, `/admin/rooms` | Номерной фонд. |
| GET, PATCH | `/api/v1/admin/bookings`, `/admin/bookings/{id}/status` | Управление бронированиями. |

Для `/incidents`: `page`, `limit`, `status`, `priority`, `hotel_id`.

## Проверка

```bash
go vet ./...
go test ./...
go build ./...
cd frontend && npm run build
```

CI запускает проверки backend и production-сборку frontend на `push` и `pull_request`.
