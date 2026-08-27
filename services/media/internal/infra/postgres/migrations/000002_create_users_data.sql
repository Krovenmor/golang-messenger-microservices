-- +goose Up
create table if not exists UsersData (
    user_id uuid references Users(user_id) on delete cascade,
    media_id uuid not null,

    is_public boolean not null,
    media_type text not null,
    media_subtype text not null,
    media_size bigint not null check (media_size > 0),

    added_at timestamptz default now(),

    primary key (user_id, media_id)
);

-- +goose Down
drop table if exists UsersData;