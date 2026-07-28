select messageid, senderid, message, createdat
from messages
where chatid = $1 and messageid >= $2
order by messageid
limit $3;