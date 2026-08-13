select
    message_id, emoji, json_agg(user_id) as users
from message_reactions
where message_id = ANY($1::uuid[])
group by message_id, emoji
order by message_id desc;