# Go Messenger Microservices
## Стек технологий

* **Язык:** Go
* **Архитектура & DI:** Event-Driven Architecture, Clean Architecture, Uber FX
* **Базы данных & Брокеры:** PostgreSQL (pgx), Redis (Streams, Pub/Sub, KV Cache)
* **Сеть & Транспорт:** REST API, WebSockets, HTTP Reverse Proxy
* **Безопасность:** JWT (Ed25519 асимметричные ключи), Rate Limiting, IP/User Blacklisting
* **Инфраструктура:** Docker, Docker Compose

---

## Архитектура и сервисы
```text
  [Client]
      |
      ↓
[edge-service] ---> [web-service] (статика фронта)
      |
      ↓
[gateway-service]
      |
      ↓
[internal-services]
      ^
      |
[redis PubSub/Streams]
```

### Описание компонентов:

- **`services/edge` (Edge / Reverse Proxy):**
  Единая внешняя точка входа для фронтенда и API-запросов. Отвечает за маршрутизацию и прямую отдачу публичной статики (аватарки пользователей).

- **`services/gateway` (API Gateway & Traffic Control):**
  Внутренний шлюз для защиты сервисов. Считает обращения по IP/User ID (Rate Limiter), блокирует нарушителей и слушает служебные события на бан от других микросервисов.

- **`services/auth` (Сервис авторизации):**
  Единственный сервис, хранящий приватный ключ для подписи асимметричных JWT-токенов. Отвечает за регистрацию и аутентификацию учетных записей.

- **`services/profile` (Сервис профилей):**
  Хранит пользовательские профили, никнеймы и метаданные. При создании профиля публикует событие в Redis Streams для синхронизации с другими доменами.

- **`services/msg` (Сервис сообщений):**
  Управляет созданием диалогов/групповых чатов, сохранением истории переписки и реакциями.

- **`services/media` (Медиа-хранилище):**
  Принимает пользовательские файлы, ведет строгий учет дисковых квот (занятое/доступное место) и валидирует входные форматы.

- **`services/status` (Сервис статусов):**
  Отслеживает сетевой статус пользователей (онлайн/офлайн) через временные ключи и события в Redis.

- **`services/ws` (WebSocket Gateway):**
  Управляет пулом постоянных соединений с клиентами. Подписывается на шину событий Redis (Pub/Sub) и транслирует входящие сообщения и системные уведомления в браузеры.

- **`services/email` (Email Worker):**
  Фоновый воркер для отправки служебных писем, читающий задачи из очереди Redis.

---

## Структура проекта
```text
├── config/       # Конфигурационные файлы для микросервисов
├── data/         # Директория для локального хранения статики и медиа
├── deploy/       # Dockerfile и docker-compose.yml для развертывания
├── pkg/          # Переиспользуемые пакеты (redis, jwt-auth, config, broker-events)
├── services/     # Исходный код микросервисов
├── web/          # Клиентская часть для демонстрации работы API
└ api.md          # Полная документация по API
```

## Запуск проекта
### Ручной запуск
1. Заполните учетные данные SMTP в файле `config/email.env` для работы `email-service`.
2. Сгенерируйте пару ключей Ed25519 для JWT и поместите их в `certs/jwt/` (private.key && public.key).
```bash
mkdir -p certs/jwt
openssl genpkey -algorithm ed25519 -out certs/jwt/private.key
openssl pkey -in certs/jwt/private.key -pubout -out certs/jwt/public.key
```
3. Сгенерируйте локальный ssl сертификат и поместите в certs/ssl/ (key.pem && cert.pem)
```bash
mkdir -p certs/ssl
openssl req -x509 -newkey rsa:2048 -nodes -keyout certs/ssl/key.pem -out certs/ssl/cert.pem -days 365 -subj "/CN=localhost"
```
4. Запуск:
```bash
docker compose -f deploy/docker-compose.yml up --build -d
```