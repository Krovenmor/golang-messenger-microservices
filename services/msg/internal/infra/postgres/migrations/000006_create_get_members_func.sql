-- +goose Up
create or replace function get_members(chat_id_ref uuid)
returns json
language sql
as $$
    select
        coalesce(
            json_agg(
                json_build_object(
                    'userId', c.user_id,
                    'name', p.name,
                    'joinedAt', c.joined_at
                )
            ),
            '[]'::json
        )
    from chatmembers c
    join profiles p on p.user_id = c.user_id
    where c.chat_id = chat_id_ref
$$;

-- +goose Down
drop function if exists get_members(chat_id_ref uuid);