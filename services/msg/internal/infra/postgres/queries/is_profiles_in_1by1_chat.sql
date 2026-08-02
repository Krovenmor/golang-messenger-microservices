select chat_id
from chatmembers
where user_id in ($1, $2)
group by chat_id
having count(distinct user_id) = 2;