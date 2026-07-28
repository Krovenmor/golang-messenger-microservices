-- +goose Up
create table if not exists Profiles(
    userId uuid primary key,
    username text unique not null,
    name text not null,
    pubKey text not null,
    prvKey text not null,
    kdfSalt text not null,
    createdAt timestamp default now()
);

-- +goose Down
drop table if exists Profiles;