-- +goose Up
create index if not exists idx_messages_pair_ids
on messages(chat_id, message_id asc);

-- +goose Down
drop index if exists idx_messages_pair_ids;