update messages
set message = $4, is_redacted = true, redacted_at = now(),
    receiver_key = $5, sender_key = $6, nonce = $7
where chat_id = $1 and message_id = $2 and sender_id = $3 and is_deleted = false;