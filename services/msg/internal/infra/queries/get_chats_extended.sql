with chats as (
    select distinct
        chatid
    from chatmembers
    where userid = $1
),
chatsPeers as (
    select
        cm.chatid,
        cm.userid,
        p.name,
        p.username
    from chatmembers cm
    join profiles p on p.userid = cm.userid
    where cm.userid != $1 and cm.chatid in (select chatid from chats)
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
    cp.chatid, cp.userid, cp.name, cp.username,
    lm.messageid, lm.senderid, lm.message, lm.createdat
from chatsPeers cp
left join lastMessage lm on lm.chatid = cp.chatid
order by lm.createdat desc nulls first;