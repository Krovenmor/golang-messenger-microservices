-- +goose Up
create table if not exists Messages(
    message_id uuid primary key,
    chat_id uuid references chats(id) on delete cascade,
    sender_id uuid references profiles(user_id) on delete cascade,

    message text not null,
    sender_key text not null,
    receiver_key text not null,
    nonce text not null,

    created_at timestamptz default now(),
    is_redacted boolean default false,
    is_deleted boolean default false,
    redacted_at timestamptz default null,

    foreign key (chat_id, sender_id) references chatmembers(chat_id, user_id)
);

-- +goose Down
drop table if exists Messages;