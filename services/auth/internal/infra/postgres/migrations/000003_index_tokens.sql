-- +goose Up
create index if not exists idx_tokens_user_token on Tokens(user_id, r_token);

-- +goose Down
drop index if exists idx_tokens_user_token;
