update messages
set message = $4, redacted_at = now(),
    receiver_key = $5, sender_key = $6, nonce = $7
where chat_id = $1 and message_id = $2 and sender_id = $3 and deleted_at is null;