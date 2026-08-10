select
    u.login,
    coalesce(json_agg(t.r_token), '[]') as tokens
from users u
join tokens t on t.user_id = u.user_id
where u.user_id = $1
group by u.login;