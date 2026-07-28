select
    userid,
    joinedat
from chatmembers
where chatid = $1;