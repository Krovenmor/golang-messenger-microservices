update profiles
set additional = jsonb_set(
    additional,
    '{avatars}',
    coalesce(additional->'avatars', '[]'::jsonb) || to_jsonb($2::text)
)
where user_id = $1;