select  userId, name, pubKey, prvKey, kdfSalt
from profiles
where userid = $1;