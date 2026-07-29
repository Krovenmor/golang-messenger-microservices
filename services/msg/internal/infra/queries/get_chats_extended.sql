with chats as (
    select distinct
        chatid
    from chatmembers
    where userid = $1
),
lastMessage as (
    select distinct on (chatid)
        chatid,
        senderid,
        message,
        messageid,
        createdat
    from messages
    where chatid in (select chatid from chats)
    order by chatid, messageid desc
)
select
    c.chatid, lm.messageid, lm.senderid, lm.message, lm.createdat
from chats c
left join lastMessage lm on lm.chatid = c.chatid
order by lm.createdat desc nulls first;