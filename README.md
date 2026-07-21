# Пульс — приложение для опросов

## Назначение

Пульс позволяет создать многовопросный опрос, поделиться публичной ссылкой и собрать по одному варианту ответа на каждый вопрос. Автор получает административный токен и ссылку на результаты.

Возможности:

- создание опроса с названием, вопросами и вариантами ответа;
- UUID для опросов, вопросов и вариантов;
- публичный просмотр без результатов голосования;
- один отправленный голос на опрос в браузере благодаря HttpOnly-cookie;
- проверка, что выбранный вариант относится к соответствующему вопросу и опросу;
- результаты с количеством голосов и процентом для каждого варианта;
- Vue SPA с экранами создания, голосования и результатов;
- PostgreSQL-миграция, запускаемая сервером при старте.

## Стек

- Go 1.21, `net/http`, chi v5, `log/slog`;
- PostgreSQL 16, pgx/v5 и pgcrypto для UUID по умолчанию;
- Vue 3, Vue Router 4, Vite 6, Tailwind CSS;
- Docker Compose и трёхэтапный Dockerfile: Node 20, Go 1.21, Alpine 3.19.

## Архитектура

```text
cmd/server/main.go                 запуск, конфигурация, миграция, static/SPA
internal/domain/models.go          доменные модели и JSON-модели
internal/repository/interfaces.go  контракт хранилища и общие ошибки
internal/repository/postgres/      SQL через pgx: polls.go, votes.go
internal/service/poll_service.go   валидация, токен администратора, правила голосования
internal/handler/poll_handler.go   HTTP-маршруты, JSON, cookies и статусы
migrations/001_init.sql            схема polls/questions/options/votes
frontend/src/                       Vue-приложение и его маршруты
```

Handler зависит только от service и общих repository-интерфейсов; прямого подключения PostgreSQL в нём нет. Создание опроса вставляет все сущности в одной транзакции.

## Требования и запуск

Нужны Docker с Compose plugin. Самый простой запуск из корня проекта:

```sh
docker compose up --build
```

После запуска приложение доступно на <http://localhost:8080>. Compose по умолчанию использует PostgreSQL:

```text
DB_HOST=db       DB_PORT=5432       DB_NAME=polls
DB_USER=postgres DB_PASSWORD=postgres DB_SSLMODE=disable
PORT=8080
```

Файл `.env` необязателен для значений по умолчанию Compose. Для своих параметров скопируйте пример и измените его:

```sh
cp .env.example .env
docker compose up --build
```

Сервер также принимает `DATABASE_URL`; он имеет приоритет над `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` и `DB_SSLMODE`. В образе миграция читается из `migrations/001_init.sql`, поэтому запускать сервер нужно из корня либо использовать Docker-образ.

## Локальная разработка

Нужны Go 1.21+, Node.js 20+, npm и доступный PostgreSQL. Запустите базу отдельно, затем:

```sh
export DATABASE_URL='postgres://postgres:postgres@localhost:5432/polls?sslmode=disable'
go run ./cmd/server
```

В другом терминале для frontend:

```sh
cd frontend
npm ci
npm run dev
```

Проверки проекта:

```sh
go test ./...
go build ./cmd/server
cd frontend && npm run build
docker compose config
docker build .
```

## Миграции и данные

`migrations/001_init.sql` идемпотентно создаёт расширение `pgcrypto`, таблицы `polls`, `questions`, `options`, `votes` и индекс голосов. Сервер выполняет этот файл после успешного подключения к БД при каждом старте. Отдельного migration runner и версионирования следующих файлов сейчас нет. Данные PostgreSQL хранятся в volume `postgres-data`.

## API

Все API-маршруты имеют префикс `/api`. JSON использует `snake_case`. Ошибки всегда имеют форму:

```json
{"error":"readable message"}
```

### Создание опроса

`POST /api/polls`, `Content-Type: application/json`

Запрос:

```json
{
  "title": "Планы на лето",
  "questions": [
    {"text": "Куда поедем?", "options": ["Море", "Горы"]},
    {"text": "Когда?", "options": ["Июнь", "Август"]}
  ]
}
```

Успех: `201 Created`.

```json
{
  "id": "8f2c...",
  "admin_token": "...",
  "public_link": "/polls/8f2c...",
  "admin_link": "/polls/8f2c.../results?admin_token=..."
}
```

Невалидный JSON или данные: `400` и `{"error":"invalid request"}`. Ограничения: название 1–500 символов, 1–50 вопросов, текст вопроса 1–1000 символов, 2–50 вариантов на вопрос, вариант до 500 символов; варианты после trim не должны повторяться без учёта регистра.

### Получение опроса

`GET /api/polls/{id}` возвращает `200 OK`:

```json
{
  "id": "8f2c...",
  "title": "Планы на лето",
  "questions": [
    {"id": "...", "text": "Куда поедем?", "options": [
      {"id": "...", "text": "Море"},
      {"id": "...", "text": "Горы"}
    ]}
  ]
}
```

Этот ответ не содержит количества голосов. Пустой или неизвестный id даёт `400`/`404`; ошибка БД — `500`.

### Голосование

`POST /api/polls/{id}/vote`

```json
{
  "answers": [
    {"question_id": "...", "option_id": "..."},
    {"question_id": "...", "option_id": "..."}
  ]
}
```

Нужно передать ровно один ответ для каждого вопроса; повтор вопроса запрещён. Сервер дополнительно проверяет принадлежность вопроса и варианта этому опросу. Успех: `201 Created`, `{"message":"vote recorded"}`.

После успеха сервер устанавливает cookie `voted_{poll_id}=true` с `Path=/`, `HttpOnly`, `SameSite=Lax` и сроком около одного года. Если такая cookie уже присутствует, запрос отклоняется до записи в БД: `409 Conflict`, `{"error":"already voted"}`. Невалидный запрос даёт `400`, неизвестный опрос — `404`, прочая ошибка — `500`.

Cookie — практическая защита от повторной отправки в одном браузере, а не учёт личности: её можно удалить или обойти другим клиентом. В базе нет voter identity и глобального ограничения на одного человека.

### Результаты

`GET /api/polls/{id}/results` доступен публично без query-параметров. Если передан `admin_token`, он проверяется:

`GET /api/polls/{id}/results?admin_token=...`

Ответ `200 OK`:

```json
{
  "id": "8f2c...",
  "title": "Планы на лето",
  "questions": [
    {"id": "...", "text": "Куда поедем?", "options": [
      {"id": "...", "text": "Море", "votes": 3, "percentage": 75},
      {"id": "...", "text": "Горы", "votes": 1, "percentage": 25}
    ]}
  ]
}
```

Неверный переданный токен: `403` и `{"error":"forbidden"}`. Неизвестный опрос: `404`; прочая ошибка: `500`. Проценты считаются отдельно внутри каждого вопроса; при нулевых голосах равны `0`.

`GET /health` возвращает `200` без JSON-тела.

## Frontend-маршруты

- `/` — главная и созданные в этом браузере опросы;
- `/create` — создание опроса;
- `/poll/:id` и `/polls/:id` — публичное голосование;
- `/poll/:id/results` и `/polls/:id/results` — результаты.

Frontend сохраняет созданные опросы и admin token в локальном storage браузера. Backend для неизвестных не-API путей отдаёт SPA `index.html`; пути `/api` не маскируются frontend fallback.
