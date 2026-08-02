select
    user_id,
    joined_at
from chatmembers
where chat_id = $1;