select
    ch.created_at,
    get_members(c.chat_id) AS members_json
from chatmembers c
join chats ch on ch.id = c.chat_id
where chat_id = $1
group by c.chat_id, ch.created_at;