select message_id, sender_id,
    message, sender_key, receiver_key, nonce,
    created_at, redacted_at, deleted_at, reply_to_id
from messages
where chat_id = $1 and message_id < $2
order by message_id desc
limit $3;