update messages
set message = $4, is_redacted = true, redacted_at = now()
where chat_id = $1 and message_id = $2 and sender_id = $3 and is_deleted = false;