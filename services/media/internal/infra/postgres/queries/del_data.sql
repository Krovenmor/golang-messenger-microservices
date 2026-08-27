with del_data as (
    delete from usersdata
    where user_id = $1 and media_id = $2
    returning user_id, media_size
)
update users u
set space_filled = u.space_filled - d.media_size,
files_count = u.files_count - 1
from del_data d
where u.user_id = d.user_id;