select 1 as is
from chatmembers
where chatid = $1 and userid = $2;