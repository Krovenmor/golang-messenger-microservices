-- +goose Up
create table if not exists ChatMembers (
    chatId uuid references chats(id) on delete cascade,
    userId uuid references profiles(userid) on delete cascade,
    joinedAt timestamp default now(),

    primary key (chatId, userId)
);

-- +goose Down
drop table if exists ChatMembers;