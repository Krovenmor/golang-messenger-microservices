insert into message_reactions(message_id, user_id, emoji)
select
    $1, $2, $3
where exists(
    select 1
    from messages
    where message_id = $1 and chat_id = $4
);