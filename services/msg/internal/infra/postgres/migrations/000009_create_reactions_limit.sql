-- +goose Up
-- +goose StatementBegin
create or replace function check_reactions_count()
returns trigger
language plpgsql
as $$
declare
    counter int;
begin
    select count(*)
    into counter
    from message_reactions
    where message_id = NEW.message_id and user_id = NEW.user_id;

    if counter >= 3 then
        raise exception 'too_many_reactions';
    end if;

    return NEW;
end;
$$;
-- +goose StatementEnd

create trigger trg_limit_user_reactions
before insert on message_reactions
for each row
execute function check_reactions_count();

-- +goose Down
drop trigger if exists trg_limit_user_reactions on message_reactions;
drop function if exists check_reactions_count();