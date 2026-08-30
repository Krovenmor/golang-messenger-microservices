update profiles
set additional = jsonb_set(
    additional,
    '{bio}',
    to_jsonb($2::text)
)
where user_id = $1;