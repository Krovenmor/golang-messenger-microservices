select
    max_space - space_filled as available_space
from users
where user_id = $1;