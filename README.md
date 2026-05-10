# go_core_rest_api

Библиотека на Go для запуска REST API поверх PostgreSQL: маршруты и вызовы хранимых функций задаются через YAML, конфигурация сервера и БД — через JSON.

## Возможности

- HTTP-сервер на [Fiber](https://github.com/gofiber/fiber): группы маршрутов, GET/POST/PUT/DELETE.
- Вызов функций БД через `SELECT имя_функции()` или `SELECT имя_функции($1)` с телом запроса как JSON.
- Валидация тела запроса по описанию полей в YAML ([go-playground/validator](https://github.com/go-playground/validator)).
- Пул соединений к PostgreSQL с разумными значениями по умолчанию и настройкой через конфиг.
- Запросы к БД через `QueryRowContext`: отмена при разрыве клиентского соединения (контекст Fiber).
- Произвольные маршруты через `InitCustomHandler` (после инициализации сервиса).

## Требования

- Go **1.25.4** (см. `go.mod`).
- PostgreSQL и драйвер `postgres` (`github.com/lib/pq`).

## Установка

```bash
go get github.com/TelitsynNikita/go_core_rest_api
```

Импорт:

```go
import "github.com/TelitsynNikita/go_core_rest_api"
```

## Быстрый старт

Процесс читает файлы относительно **текущей рабочей директории** процесса:

- `./configs/config.json` — сервер и БД;
- `./apis.yaml` — описание API.

```bash
cd example
go run .
```

Скопируйте или адаптируйте `example/configs/config.json` и `example/apis.yaml` под свой проект.

### Жизненный цикл сервиса

1. `NewApiService()` — пустой сервис.
2. `InitApiService()` — загрузка JSON-конфига, YAML API, подключение к БД, создание Fiber-приложения и регистрация маршрутов из YAML (для `type: postgresql`).
3. При необходимости — один или несколько вызовов `InitCustomHandler(...)` для своих обработчиков.
4. `RunService()` — прослушивание порта из `server_config.port`.

### Конфигурация (`configs/config.json`)

| Поле | Описание |
|------|-----------|
| `server_config.port` | Порт HTTP-сервера |
| `db_config` | Параметры подключения к PostgreSQL (`user`, `password`, `host`, `port`, `db_name`, `driver`, `ssl_mode` и др.) |
| `db_config.max_open_conns` | Максимум открытых соединений (по умолчанию **25**) |
| `db_config.max_idle_conns` | Простаивающих соединений в пуле (по умолчанию **5**, не больше `max_open_conns`) |
| `db_config.conn_max_lifetime_sec` | Максимальное время жизни соединения в секундах (по умолчанию **300**) |
| `db_config.conn_max_idle_time_sec` | Как долго простаивающее соединение держится в пуле (по умолчанию **900**) |

Нулевые или отсутствующие поля пула подставляют значения по умолчанию из кода (`applyPoolSettings` в `init_db.go`).

### Описание API (`apis.yaml`)

Файл — вложенная карта: **имя группы** (префикс пути) → **относительный путь** → настройки эндпоинта.

| Поле в YAML | Описание |
|-------------|-----------|
| `type` | `postgresql` — маршрут ведёт в `SelectFunction` с полем `call`. Значение `custom` в YAML зарезервировано; такие записи **не** регистрируют обработчик автоматически — используйте `InitCustomHandler`. |
| `method` | HTTP-метод: `get`, `post`, `put`, `delete`. |
| `call` | Имя функции в БД (для `postgresql`), например `schema.function_name`. |
| `body` | Описание полей JSON-тела для валидации: для каждого поля задаются `type` и `tag` (правила validator). Пустой `body` допустим. |
| `is_slice` | Если `true`, тело ожидается как JSON-массив объектов согласно `body`. |

Итоговый путь для эндпоинта: `/{группа}/{url}`, например группа `api` и ключ `get-all-cs-equipment` дают маршрут `GET /api/get-all-cs-equipment` (точный URL зависит от объявления в YAML).

### Пользовательские маршруты (`InitCustomHandler`)

После `InitApiService()` можно зарегистрировать свой обработчик Fiber. Разрешены методы: GET, POST, PUT, DELETE (в любом регистре); иной метод приведёт к `panic`.

Сигнатура:

```go
func (a *ApiService) InitCustomHandler(group, url, method string, fn func(c *fiber.Ctx) error)
```

Маршрут: `METHOD /{group}/{url}` (метод нормализуется в верхний регистр).

Для прямых запросов к пулу используйте `apiService.DB.SQLX` (`*sqlx.DB`) и по возможности контекст из `c.UserContext()`. Универсальный разбор ответа функции, возвращающей JSON, даёт `Database.SelectFunction`.

Пример см. в `example/server.go`.

## Структура репозитория

| Файл | Назначение |
|------|------------|
| `core.go` | Инициализация конфигурации, YAML API, БД и Fiber-приложения |
| `init_db.go` | Подключение к БД, пул, `SelectFunction`, поле `Database.SQLX` |
| `init_handlers.go` | Регистрация хэндлеров из YAML и вызов функций БД |
| `config_reader.go` | Чтение JSON-конфига |
| `yaml_reader.go` | Чтение и валидация описания API из YAML |
| `example/` | Пример запуска сервиса и кастомного маршрута |

## Слой базы данных

- Запросы к функциям без лишней транзакции: один `SELECT` через пул, без `Begin`/`Commit` на каждый вызов.
- `SelectFunction(ctx, имяФункции, body)` принимает контекст; из стандартных HTTP-хэндлеров передаётся `c.UserContext()` из Fiber.
- Для произвольного SQL или обхода обёртки используйте `DB.SQLX` с `QueryRowContext` / `QueryContext` и т.д.

Подробности изменений — в [CHANGELOG.md](CHANGELOG.md). Идеи развития — в [FUTURE_FEATURES.md](FUTURE_FEATURES.md).

## Лицензия

Уточните при необходимости у автора репозитория.
