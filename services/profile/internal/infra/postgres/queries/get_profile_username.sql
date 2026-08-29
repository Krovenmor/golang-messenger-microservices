select 
    user_id, name, user_name, pub_key, prv_key, kdf_salt, key_nonce, created_at, additional
from profiles
where user_name = $1;