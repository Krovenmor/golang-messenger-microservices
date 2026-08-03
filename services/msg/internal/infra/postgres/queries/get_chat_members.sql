select
    c.user_id,
    p.name,
    c.joined_at
from chatmembers c
join profiles p on p.user_id = c.user_id
where chat_id = $1