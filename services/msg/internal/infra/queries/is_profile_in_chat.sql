select 1 as is
from chatmembers
where chat_id = $1 and user_id = $2;