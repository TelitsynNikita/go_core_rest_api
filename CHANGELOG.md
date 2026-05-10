# Changelog

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.0.0/).

## [Unreleased]

## [0.1.0] — 2026-05-10

### Added

- Документация: `README.md`, этот файл, `FUTURE_FEATURES.md`.
- Настройка пула соединений PostgreSQL в `DBConfig`: `max_open_conns`, `max_idle_conns`, `conn_max_lifetime_sec`, `conn_max_idle_time_sec` с значениями по умолчанию в `applyPoolSettings` (`init_db.go`).
- Поддержка контекста в `SelectFunction(ctx context.Context, ...)`: выполнение через `QueryRowContext`; в хэндлерах передаётся `c.UserContext()` из Fiber (`init_handlers.go`).

### Changed

- Удалена лишняя транзакция вокруг одного `SELECT`: запросы выполняются напрямую через пул без `Begin`/`Commit`, что убирает двойное использование соединений на один вызов.

### Fixed

- Поведение при вызовах API к БД без «лишней» транзакции и с ограничением пула снижает риск исчерпания подключений при нагрузке.
