select chatid
from chatmembers
where userid in ($1, $2)
group by chatid
having count(distinct userid) = 2;