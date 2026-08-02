-- +goose Up
create table if not exists Profiles(
    user_id uuid primary key,
    user_name text unique not null,
    name text not null,
    pub_key text not null,
    prv_key text not null,
    kdf_salt text not null,
    created_at timestamp default now()
);

-- +goose Down
drop table if exists Profiles;