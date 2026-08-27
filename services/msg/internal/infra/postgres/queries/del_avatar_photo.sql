update profiles
set additional = jsonb_set(
    additional,
    '{avatars}',
    (additional->'avatars') - $2::text
)
where user_id = $1;