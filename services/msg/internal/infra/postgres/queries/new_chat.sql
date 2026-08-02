WITH new_chat AS (
    INSERT INTO chats (id)
    VALUES ($1)
    RETURNING id
)
INSERT INTO chatmembers (chat_id, user_id)
VALUES
    ($1, $2),
    ($1, $3);