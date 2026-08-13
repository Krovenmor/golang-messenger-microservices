-- +goose Up
create table if not exists message_reactions (
    message_id uuid not null references messages(message_id) on delete cascade,
    user_id uuid not null references profiles(user_id) on delete cascade,
    emoji varchar(64) not null references emojis(emoji) on delete restrict,

    created_at timestamptz not null default now(),

    PRIMARY KEY (message_id, user_id, emoji)
);

-- +goose Down
drop table if exists message_reactions;