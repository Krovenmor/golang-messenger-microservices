-- +goose Up
alter table profiles add constraint check_avatars_in_additional
check (
    coalesce(jsonb_array_length(additional->'avatars'), 0) <= 10
);

-- +goose Down
alter table profiles drop constraint check_avatars_in_additional;