-- +goose Up
create table if not exists Chats(
    id uuid primary key,
    created_at timestamp default now()
);

-- +goose Down
drop table if exists Chats;