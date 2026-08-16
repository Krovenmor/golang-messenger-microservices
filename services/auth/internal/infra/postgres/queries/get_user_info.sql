select
    login,
    email
from users
where user_id = $1;