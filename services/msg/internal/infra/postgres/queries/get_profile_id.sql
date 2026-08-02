select 
    user_id, name, user_name, pub_key, prv_key, kdf_salt, created_at
from profiles
where user_id = $1;