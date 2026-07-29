select 1
from tokens
where user_id = $1 and r_token = $2;