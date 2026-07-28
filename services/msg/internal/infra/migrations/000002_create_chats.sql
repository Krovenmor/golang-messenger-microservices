-- +goose Up
create table if not exists Chats(
    id uuid primary key,
    fPerson uuid references profiles(userid) on delete cascade,
    sPerson uuid references profiles(userid) on delete cascade,

    constraint unique_pair unique(fPerson, sPerson)
);

-- +goose Down
drop table if exists Chats;