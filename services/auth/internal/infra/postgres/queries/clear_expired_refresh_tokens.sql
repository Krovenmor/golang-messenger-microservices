delete from tokens
where user_id = $1 and exp_at < now();