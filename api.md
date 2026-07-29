**Auth API:**  
-

- "POST /api/auth/register" - register a new user  
Income:
{
    "login": "YourLogin",
    "password": "YourPassword"
}  

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
    "Name": "Your Public Name",
    "UserName": "Your Unique UserName",  
    "PubKey": "Public key",  
    "EncryptedPrvKey": "Encrypted Private Key",  
    "KDFSalt": "KDF Salt"  
}  

- "GET /api/msg/profile" - get yours profile
Outcome:  
{  
    "Name": "Name",
    "UserName": "User Name",  
    "PubKey": "Public key",  
    "EncryptedPrvKey": "Encrypted Private Key",  
    "KDFSalt": "KDF Salt"  
    "CreatedAt": "time, example: '2026-07-28T19:56:51.855208Z'"
}  

- "GET /api/msg/profile/{target}" - get other profile  
{target} - UserId (UUID) or UserName  
Outcome:  
{  
    "UserId": "UUID of profile",  
    "Name": "Name",  
    "UserName": "User Name",  
    "PubKey": "Public key"  
}  

- "GET /api/msg/profile/chats" - get all your chats  
Outcome:  
[  
    "UUID1", "UUID2", ...  
]  

- "GET /api/msg/profile/chats/full" - get all your chats with extended info  
Outcome:
[
    {
        "ChatId": "UUID",
        "MessageId": "Last Message Id in chat, UUIDv7",
        "SenderId": "Who wrote that message, UUID",
        "Message": "text ...",
        "CreatedAt": "time UTC"
    }, ...
]

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

- "GET /api/msg/chat/{uuid}/info" - get chat info  
Outcome:  
{  
    "CreatedAt": "Time when it was created",
    "Members": [
        {"UserId": "UUID", "JoinedAt": "Time"}, {...} ...
    ]
}