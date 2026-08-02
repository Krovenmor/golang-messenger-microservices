select message_id, sender_id, message, created_at,
    is_redacted, is_deleted, redacted_at
from messages
where chat_id = $1 and message_id < $2
order by message_id desc
limit $3;