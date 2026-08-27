select
    max_space, space_filled, files_count
from users
where user_id = $1;