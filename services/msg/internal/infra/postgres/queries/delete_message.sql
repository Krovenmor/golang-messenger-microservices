update messages
set message = '', sender_key = '', receiver_key = '', nonce = '', reply_to_id = null, deleted_at = now()
where chat_id = $1 and message_id = $2 and sender_id = $3 and deleted_at is null;