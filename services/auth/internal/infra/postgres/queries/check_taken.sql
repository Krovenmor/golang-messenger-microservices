select
    exists(select 1 from users where login = $1),
    exists(select 1 from users where lower(email) = lower($2));
