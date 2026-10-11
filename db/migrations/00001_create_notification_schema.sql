-- Tables are added after ADR-03 (storage and secrets) and ADR-04 (state model)
-- are accepted. See docs/Notitek - CẤU TRÚC MÃ NGUỒN.md §7 for table ownership.

-- +goose Up
CREATE SCHEMA IF NOT EXISTS notification;

-- +goose Down
DROP SCHEMA IF EXISTS notification;
