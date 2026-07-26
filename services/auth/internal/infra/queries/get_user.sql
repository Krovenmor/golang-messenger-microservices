select u.userid
from users u
where u.login = $1 and u.password = $2;