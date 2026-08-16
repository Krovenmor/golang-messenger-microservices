-- +goose Up
create table if not exists Users (
    user_id uuid primary key,
    login VARCHAR(50) unique not null,
    password VARCHAR(128) not null,
    email varchar(255) not null,

    created_at timestamptz default now()
);

create unique index idx_unique_email on Users(lower(email));

-- +goose Down
drop table if exists Users;