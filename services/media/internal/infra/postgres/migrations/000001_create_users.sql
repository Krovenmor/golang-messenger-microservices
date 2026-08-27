-- +goose Up
create table if not exists Users (
    user_id uuid primary key,
    max_space bigint not null default 104857600,
    space_filled bigint not null default 0 check (space_filled <= max_space),
    files_count int default 0 check (files_count >= 0),

    created_at timestamptz default now()
);

-- +goose Down
drop table if exists Users;