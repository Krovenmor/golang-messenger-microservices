with chats_ids as (
    select distinct
        chat_id
    from chatmembers
    where user_id = $1
),
lastMessage as (
    select distinct on (chat_id)
        chat_id,
        created_at,
        json_build_object(
            'messageId', message_id,
            'senderId', sender_id,
            'message', message,
            'senderKey', sender_key,
            'receiverKey', receiver_key,
            'nonce', nonce,
            'createdAt', created_at,
            'redactedAt', redacted_at,
            'deletedAt', deleted_at,
            'replyToId', reply_to_id
        ) as last_message_json
    from messages
    where chat_id in (select chat_id from chats_ids)
    order by chat_id, message_id desc
),
members_aggregated AS (
    select
        c.chat_id,
        get_members(c.chat_id) AS members_json
    from chatmembers c
    join profiles p on p.user_id = c.user_id
    where chat_id in (select chat_id from chats_ids)
    group by c.chat_id
)
select
    -- Chat
    c.chat_id, ch.created_at,
    -- Members
    coalesce(ma.members_json, '[]'::json) as chat_members,
    -- Last Message Info
    lm.last_message_json as last_message
from chats_ids c
join chats ch on ch.id = c.chat_id
left join members_aggregated ma on ma.chat_id = c.chat_id
left join lastMessage lm on lm.chat_id = c.chat_id
order by lm.created_at desc nulls first;