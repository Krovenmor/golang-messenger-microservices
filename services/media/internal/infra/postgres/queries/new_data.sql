with upsert_user as (
    insert into users (user_id, space_filled, files_count)
    values ($1, $5, 1)
    on conflict (user_id)
    do update
    set space_filled = users.space_filled + excluded.space_filled,
    files_count = users.files_count + excluded.files_count
    returning user_id
)
insert into usersdata (user_id, media_id, media_type, media_subtype, media_size, is_public)
select
    u.user_id, $2, $3, $4, $5, $6
from upsert_user u;