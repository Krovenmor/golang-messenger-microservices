-- +goose Up
create table if not exists emojis (
    emoji varchar(64) primary key,
    created_at timestamptz default now()
);

-- +goose Down
drop table if exists emojis;