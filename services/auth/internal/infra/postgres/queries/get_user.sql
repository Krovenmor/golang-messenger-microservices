select u.user_id, u.password
from users u
where u.login = $1;