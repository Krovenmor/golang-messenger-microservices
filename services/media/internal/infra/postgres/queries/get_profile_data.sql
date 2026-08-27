select
    media_id, media_type, media_subtype, media_size, is_public, added_at
from usersdata
where user_id = $1 and media_id < $2
order by media_id desc
limit $3;