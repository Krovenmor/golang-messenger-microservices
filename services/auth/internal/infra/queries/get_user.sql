select u.user_id
from users u
where u.login = $1 and u.password = $2;