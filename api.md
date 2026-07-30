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

MessageOutcomeBody:  
{  
    "MessageId": "id UUIDv7",  
    "SenderId": "User UUID",  
    "Message": "Message text",  
    "CreatedAt": "Time when it was added to server UTC, format RFC 3339, example: '2026-07-28T00:21:09.904497Z'",  
    "IsRedacted": true/false  
    "IsDeleted": true/false (if it's true message will be null)  
    "RedactedAt": "time ..." (if it wasn't redacted (not deleted, not redacted) then null)  
}

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
        "Name": "text ..."
        "Message": "text ...",
        "CreatedAt": "time UTC"
    }, ...
]  
"Name", "Message", "CreatedAt" can be null  

- "POST /api/msg/chat/new" - create a new chat  
Income:  
{  
    "UserId": "User UUID to create with"  
}  
Outcome:  
{  
    "ChatId": "Chat UUID"  
}  

- "POST /api/msg/chat/{ChatId}" - post a message to a chat  
Income:  
{  
    "Message": "Message text"  
}  
Outcome:  
{  
    "MessageId": "MessageId UUIDv7"  
}

- "GET /api/msg/chat/{ChatId}" - get chat history  
Query params:  
'from' - messageId from (UUIDv7) (if you don't provide 'from', you'll get last 'q' messages)
'q' - quantity (int) (if you don't provide 'q', q will be default = 10)
Outcome:  
[  
    MessageOutcomeBody, ...  
]  

- "GET /api/msg/chat/{ChatId}/info" - get chat info  
Outcome:  
{  
    "CreatedAt": "Time when it was created",
    "Members": [
        {"UserId": "UUID", "JoinedAt": "Time"}, {...} ...
    ]
}

- "GET /api/msg/chat/{ChatId}/message/{MessageId}" - get message info  
Outcome:  
{ MessageOutcomeBody }  

- "PUT /api/msg/chat/{ChatId}/message/{MessageId}" - change message text  
Income:  
{  
    "Message": "Message text"  
}

- "DELETE /api/msg/chat/{ChatId}/message/{MessageId}" - delete message from chat  