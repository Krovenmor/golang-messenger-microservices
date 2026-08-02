update messages
set message = '', is_deleted = true, redacted_at = now()
where chat_id = $1 and message_id = $2 and sender_id = $3 and is_deleted = false;