**Auth API:**  
-

- **"POST /api/auth/register"** - register a new user  
Income:
```json
{
    "login": "YourLogin",
    "password": "YourPassword"
}  
```

- **"POST /api/auth/login"** - LogIn, get JWT tokens  
Income:
```json
{
    "login": "YourLogin",
    "password": "YourPassword"
}  
```
Outcome:  
```json
{
    "accessToken": "YourAccessToken",
    "refreshToken": "YourRefreshToken"
}  
```

- **"POST /api/auth/update"** - Update your tokens  
Income:
```json
{
    "refreshToken": "YourRefreshToken"
}  
```
Outcome:  
```json
{
    "accessToken": "YourNewAccessToken",
    "refreshToken": "YourNewRefreshToken"
}  
```

**Message API:**  
-
**All endpoints works only with JWT token from Auth API**  
**Usage: Authorization -> Bearer Token -> Access token from AuthApi**  

MessageOutcomeBody:  
```json
{  
    "MessageId": "id UUIDv7",  
    "SenderId": "User UUID",  
    "Message": "Message text",  
    "CreatedAt": "Time when it was added to server UTC, format RFC 3339, example: '2026-07-28T00:21:09.904497Z'",  
    "IsRedacted": true/false,
    "IsDeleted": true/false (if it's true message will be null),
    "RedactedAt": "time ..." (if it wasn't redacted (not deleted, not redacted) then null)
}
```

- **"POST /api/msg/profile/new"** - creates a new profile  
Income:  
```json
{  
    "Name": "Your Public Name",
    "UserName": "Your Unique UserName",  
    "PubKey": "Public key",  
    "EncryptedPrvKey": "Encrypted Private Key",  
    "KDFSalt": "KDF Salt"  
}  
```

- **"GET /api/msg/profile"** - get yours profile
Outcome:  
```json
{  
    "UserId": "UUID of profile",
    "Name": "Name",
    "UserName": "User Name",  
    "PubKey": "Public key",  
    "EncryptedPrvKey": "Encrypted Private Key",  
    "KDFSalt": "KDF Salt",
    "CreatedAt": "time, example: '2026-07-28T19:56:51.855208Z'"
}  
```

- **"GET /api/msg/profile/{target}"** - get other profile  
{target} - UserId (UUID) or UserName  
Outcome:  
```json
{  
    "UserId": "UUID of profile",  
    "Name": "Name",  
    "UserName": "User Name",  
    "PubKey": "Public key",
    "CreatedAt": "time, example: '2026-07-28T19:56:51.855208Z'"
}  
```

- **"GET /api/msg/profile/chats"** - get all your chats  
Outcome:  
```json
[  
    "UUID1", "UUID2", ...  
]  
```

- **"GET /api/msg/profile/chats/full"** - get all your chats with extended info  
Outcome:  
```json
[
    {
        "ChatId": "UUID",
        "Members": [
            {
                "UserId": "UUID", 
                "Name": "Name",
                "JoinedAt": "Time"
            }, ...
        ]
        "LastMessage": MessageOutcomeBody (can be null)
    }, ...
]  
```  

- **"POST /api/msg/chat/new"** - create a new chat  
Income:  
```json
{  
    "UserId": "User UUID to create with"  
}  
```
Outcome:  
```json
{  
    "ChatId": "Chat UUID"  
}  
```

- **"POST /api/msg/chat/{ChatId}"** - post a message to a chat  
Income:  
```json
{  
    "Message": "Message text"  
}  
```
Outcome:  
```json
{  
    "MessageId": "MessageId UUIDv7"  
}
```

- **"GET /api/msg/chat/{ChatId}"** - get chat history  
Query params:  
'from' - messageId from (UUIDv7) (if you don't provide 'from', you'll get last 'q' messages)
'q' - quantity (int) (if you don't provide 'q', q will be default = 10)
Outcome:  
```json
[  
    MessageOutcomeBody, ...  
]  
```

- **"GET /api/msg/chat/{ChatId}/info"** - get chat info  
Outcome:  
```json
{  
    "CreatedAt": "Time when it was created",
    "Members": [
        {
            "UserId": "UUID", 
            "Name": "Name",
            "JoinedAt": "Time"
        }, ...
    ]
}
```

- **"GET /api/msg/chat/{ChatId}/message/{MessageId}"** - get message info  
Outcome:  
```json
{ MessageOutcomeBody }  
```

- **"PUT /api/msg/chat/{ChatId}/message/{MessageId}"** - change message text  
Income:  
```json
{  
    "Message": "Message text"  
}
```

- **"DELETE /api/msg/chat/{ChatId}/message/{MessageId}"** - delete message from chat  

**WebSocket API:**  
-
**Endpoint works only with JWT token from Auth API**  

- **"WS: /api/ws?token={YourAccessToken}"** - start websocket connection  
By connecting to socket, you will recieve events:  
```json
{  
    "type": "EventType",  
    "payload": EventPayload  
}  
```

**Events:**  
- NewChatType:  
```json
{  
    "type": "newChat",  
    "payload": {
        "chatId": "ChatId"
    }  
}  
```
- NewMessageType:  
```json
{  
    "type": "newMessage",  
    "payload": {
        "chatId": "ChatID",
        "msgId": "MessageID"
    }  
}  
```
- MessageRedactedType:  
```json
{  
    "type": "messageRedacted",  
    "payload": {
        "chatId": "ChatID",
        "msgId": "MessageID"
    }  
}  
```
- MessageDeletedType:
```json
{  
    "type": "messageDeleted",  
    "payload": {
        "chatId": "ChatID",
        "msgId": "MessageID"
    }  
}  
```
- StatusEvent:
```json
{  
    "type": "status",  
    "payload": {
        "userId": "UUID of a user",
        "newStatus": 2,
        "eventTime": 1786069519 (unixTime)
    }  
}  
```

**Statuses:**  
- You can send your status via ws connection:  
Income: (Example)  
```json
{
    "req": 1,
    "payload": {
        "newStatus": 3
    }
}
```
Outcome: (Example)  
```json
{  
    "type": "response",  
    "payload": {
        "code": 200,
        "msg": "OK"
    }  
}  
```  
- newStatus - supports only Statuses from Status API  
- code - http codes, 200 or 400  
- All possible outcomes:
    - 200, "OK"
    - 400, "Bad JSON"
    - 400, "Bad Status"  
    - 400, "Bad UUID"
    - 400, "Bad Request"
    - 500, "Intrnal Error"

- You can also make request to track someone statuses (primary use: 1v1 chat):  
Income: (Example)  
```json
{
    "req": 2,
    "payload": {
        "userId": "{{OtherUserId}}"
    }
}
```
Outcome: (Example)  
```json
{  
    "type": "response",  
    "payload": {
        "code": 200,
        "msg": "OK"
    }  
}  
```  
- Requests:
    - 1: post your status
    - 2: track someone  

*When you connect via ws, you already get status Online, so don't send it. When you disconnect your profile automatically gain Offline status  

# **Status API:**  
## - "GET /api/status/{ProfileUUID}" - get last known profile status  
Outcome: (Example)
```json
{
    "Status": 1,
    "LastSeen": 1785881345
}
```
Statuses:  
- Offline = 1  
- Online = 2  
- Away = 3  
- Typing = 4  

LastSeen - Unix time  
If LastSeen = 0 then profile not exists or was online a very long time ago  