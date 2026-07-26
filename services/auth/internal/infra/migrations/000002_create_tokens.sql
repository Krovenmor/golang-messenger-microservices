-- +goose Up
create table if not exists Tokens (
    userId uuid references Users(userId) on delete cascade,
    rToken TEXT not null,
    expAt timestamp not null
);

-- +goose Down
drop table if exists Tokens;