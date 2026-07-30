-- +goose Up
create table if not exists ChatMembers (
    chat_id uuid references chats(id) on delete cascade,
    user_id uuid references profiles(user_id) on delete cascade,
    joined_at timestamp default now(),

    primary key (chat_id, user_id)
);

-- +goose Down
drop table if exists ChatMembers;