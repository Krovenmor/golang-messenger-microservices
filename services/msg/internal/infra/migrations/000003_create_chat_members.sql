create table if not exists chat_members (
    chatId uuid references chats(id) on delete cascade,
    userId uuid references profiles(userid) on delete cascade,
    joinedAt timestamp default now(),

    primary key (chatId, userId)
);