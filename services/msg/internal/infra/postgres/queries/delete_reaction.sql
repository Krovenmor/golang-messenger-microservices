delete
from message_reactions
where message_id = $1 and user_id = $2 and emoji = $3 and exists(
    select 1
    from messages
    where message_id = $1 and chat_id = $4
);