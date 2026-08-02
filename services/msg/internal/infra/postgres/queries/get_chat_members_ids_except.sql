select
    user_id
from chatmembers
where chat_id = $1 and user_id != $2;