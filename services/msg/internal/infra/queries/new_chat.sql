WITH new_chat AS (
    INSERT INTO chats (id)
    VALUES ($1)
    RETURNING id
)
INSERT INTO chatmembers (chatid, userid)
VALUES
    ($1, $2),
    ($1, $3);