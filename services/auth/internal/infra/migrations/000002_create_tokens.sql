-- +goose Up
create table if not exists Tokens (
    user_id uuid references Users(user_id) on delete cascade,
    r_token TEXT not null,
    exp_at timestamptz not null
);

-- +goose Down
drop table if exists Tokens;