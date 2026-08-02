with chats as (
    select distinct
        chat_id
    from chatmembers
    where user_id = $1
),
lastMessage as (
    select distinct on (chat_id)
        chat_id,
        sender_id,
        message,
        message_id,
        created_at
    from messages
    where chat_id in (select chat_id from chats)
    order by chat_id, message_id desc
)
select
    c.chat_id, lm.message_id,
    lm.sender_id, p.name,
    lm.message, lm.created_at
from chats c
left join lastMessage lm on lm.chat_id = c.chat_id
left join profiles p on p.user_id = lm.sender_id
order by lm.created_at desc nulls first;