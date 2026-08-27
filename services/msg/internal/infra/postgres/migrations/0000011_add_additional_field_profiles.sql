-- +goose Up
alter table profiles add column if not exists additional jsonb default '{}'::jsonb;

-- +goose Down
alter table profiles drop column if exists additional;