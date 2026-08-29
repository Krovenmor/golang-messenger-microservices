with upsert_user as (
    update users
    set space_filled = users.space_filled + $5,
    files_count = users.files_count + 1
    returning user_id
)
insert into usersdata (user_id, media_id, media_type, media_subtype, media_size, is_public)
select
    u.user_id, $2, $3, $4, $5, $6
from upsert_user u;