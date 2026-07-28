**Auth API:**  
-

- "POST /api/auth/register" - register a new user  
Income:
{
    "login": "YourLogin",
    "password": "YourPassword"
}  
Outcome:
Http StatusOK

- "POST /api/auth/login" - LogIn, get JWT tokens  
Income:
{
    "login": "YourLogin",
    "password": "YourPassword"
}  
Outcome:  
{
    "accessToken": "YourAccessToken",
    "refreshToken": "YourRefreshToken"
}  

- "POST /api/auth/update" - Update your tokens  
Income:
{
    "refreshToken": "YourRefreshToken"
}  
Outcome:  
{
    "accessToken": "YourNewAccessToken",
    "refreshToken": "YourNewRefreshToken"
}  

**Message API:**  
-
**All endpoints works only with JWT token from Auth API**  
**Usage: Authorization -> Bearer Token -> Access token from AuthApi**  

- "POST /api/msg/profile/new" - creates a new profile  
Income:  
{  
    "UserName": "User Name",  
    "PubKey": "Public key",  
    "EncryptedPrvKey": "Encrypted Private Key",  
    "KDFSalt": "KDF Salt"  
}  

- "GET /api/msg/profile" - get yours profile
Outcome:  
{  
    "UserName": "User Name",  
    "PubKey": "Public key",  
    "EncryptedPrvKey": "Encrypted Private Key",  
    "KDFSalt": "KDF Salt"  
}  

- "POST /api/msg/chat/new" - create a new chat  
Income:  
{  
    "UserId": "User UUID to create with"  
}  
Outcome:  
{
    "ChatId": "Chat UUID"
}

- "POST /api/msg/chat/{uuid}" - post a message to a chat  
Income:  
{  
    "Message": "Message text"  
}  
Outcome:  
{  
    "MessageId": "MessageId UUIDv7"  
}

- "GET /api/msg/chat/{uuid}" - get chat history  
Query params:  
'from' - messageId from (UUIDv7)
'q' - quantity (int)
Outcome:  
[  
    {  
    "MessageId": "id UUIDv7",  
    "SenderId": "User UUID",  
    "Message": "Message text",  
    "CreatedAt": "Time when it was added to server UTC, format RFC 3339, example: '2026-07-28T00:21:09.904497Z'"  
    }, ...  
]  