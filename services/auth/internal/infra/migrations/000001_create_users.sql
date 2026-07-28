-- +goose Up
create table if not exists Users (
    userId uuid primary key,
    login VARCHAR(50) unique not null,
    password VARCHAR(50) not null
);

-- +goose Down
drop table if exists Users;