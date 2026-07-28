-- +goose Up
create table if not exists Messages(
    messageId uuid not null,
    chatId uuid references chats(id) on delete cascade,
    sentId uuid references profiles(userid) on delete cascade,
    message text not null,
    createdAt timestamp default now(),

    primary key (MessageId, chatid)
);

-- +goose Down
drop table if exists Messages;