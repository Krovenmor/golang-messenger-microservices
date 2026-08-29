-- +goose Up
create table if not exists Profiles(
    user_id uuid primary key
);

-- +goose Down
drop table if exists Profiles;