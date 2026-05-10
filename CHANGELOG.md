# Changelog

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.0.0/).

## [Unreleased]

## [0.2.0] — 2026-05-10

### Added

- `InitCustomHandler(group, url, method, handler)` — регистрация произвольных маршрутов на Fiber после `InitApiService`; поддерживаются GET, POST, PUT, DELETE, иначе `panic`.
- Документация в README: установка модуля, жизненный цикл сервиса, таблицы по `config.json` и `apis.yaml`, раздел про кастомные маршруты и `Database.SQLX`.

### Changed

- `InitApiService` поднимает Fiber-приложение внутри себя (`initServer` сразу после `initDB`); `RunService` только вызывает `Listen`.
- Поле `Database.sqlx` заменено на экспортированное `Database.SQLX` (`*sqlx.DB`) для расширяемости и примеров с произвольными запросами.
- Внутренний `initDB` использует `a.Config.DBConfig` без отдельного аргумента конфигурации.

### Breaking

- Сигнатура `InitCustomHandler` изменена: вместо колбэка с `*Database` используется `func(c *fiber.Ctx) error`, плюс явные `group`, `url` и HTTP-метод.
- Код, обращавшийся к неэкспортированному полю `sqlx` у `Database`, нужно перевести на `SQLX`.

## [0.1.0] — 2026-05-10

### Added

- Документация: `README.md`, этот файл, `FUTURE_FEATURES.md`.
- Настройка пула соединений PostgreSQL в `DBConfig`: `max_open_conns`, `max_idle_conns`, `conn_max_lifetime_sec`, `conn_max_idle_time_sec` с значениями по умолчанию в `applyPoolSettings` (`init_db.go`).
- Поддержка контекста в `SelectFunction(ctx context.Context, ...)`: выполнение через `QueryRowContext`; в хэндлерах передаётся `c.UserContext()` из Fiber (`init_handlers.go`).

### Changed

- Удалена лишняя транзакция вокруг одного `SELECT`: запросы выполняются напрямую через пул без `Begin`/`Commit`, что убирает двойное использование соединений на один вызов.

### Fixed

- Поведение при вызовах API к БД без «лишней» транзакции и с ограничением пула снижает риск исчерпания подключений при нагрузке.
