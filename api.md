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
    "messageId": "Message UUIDv7",
    "senderId": "Sender UUIDv4",  
    // Encrypted message payload (AES-256-GCM ciphertext + Auth Tag in Base64). Evaluates to '' if deletedAt not nil
    "message": "some message...", 
    // Symmetric MessageKey encrypted with the SENDER'S public key (ECDH P-256). Allows sender to decrypt their own history. Exactly 140 Base64 characters
    "senderKey": "Some key for sender",
    // Same symmetric MessageKey encrypted with the RECEIVER'S public key.
    "receiverKey": "Some key for receiver",
    // 12-byte Initialization Vector (IV/Nonce) for AES-256-GCM encoded in Base64. Required by the client to decrypt 'message'. Exactly 16 characters
    "nonce": "a8F3b2C9m1K5p7Q0",
    // UTC timestamp when the message was persisted on the server (RFC 3339 format)
    "createdAt": "2026-07-28T00:21:09.904497Z",  
    // UTC timestamp of the last edit (RFC 3339 format). Set to null if the message was never edited
    "redactedAt": null,
    // UTC timestamp (RFC 3339 format). Set to null if the message not deleted
    "deletedAt": null,
    // Can be null
    "replyToId": "MessageId UUIDv7"
}
```

- **"POST /api/msg/profile"** - creates a new profile  
Income:  
```json
{  
    "name": "Your Public Name",
    "userName": "Your Unique UserName",  
    "pubKey": "Public key",  
    "encryptedPrvKey": "Encrypted Private Key",  
    "kdfSalt": "KDF Salt",
    "keyNonce": "Your nonce"
}  
```

- **"GET /api/msg/profile"** - get yours profile
Outcome:  
```json
{  
    "userId": "UUID of profile",
    "name": "Name",
    "userName": "User Name",  
    "pubKey": "Public key",  
    "encryptedPrvKey": "Encrypted Private Key",  
    "kdfSalt": "KDF Salt",
    "keyNonce": "Your nonce",
    "createdAt": "time, example: '2026-07-28T19:56:51.855208Z'"
}  
```

- **"GET /api/msg/profile/{target}"** - get other profile  
{target} - UserId (UUID) or UserName  
Outcome:  
```json
{  
    "userId": "UUID of profile",  
    "name": "Name",  
    "userName": "User Name",  
    "pubKey": "Public key",
    "createdAt": "time, example: '2026-07-28T19:56:51.855208Z'"
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
        "chatId": "UUID",
        "members": [
            {
                "userId": "UUID", 
                "name": "Name",
                "joinedAt": "Time"
            }, ...
        ]
        "lastMessage": MessageOutcomeBody (can be null)
    }, ...
]  
```  

- **"POST /api/msg/chat/new"** - create a new chat  
Income:  
```json
{  
    "userId": "User UUID to create with"  
}  
```
Outcome:  
```json
{  
    "chatId": "Chat UUID"  
}  
```

- **"POST /api/msg/chat/{ChatId}"** - post a message to a chat  
Income:  
```json
{   
    "message": "Message text P1",
    "senderKey": "key",
    "receiverKey": "key",
    "nonce": "nonce",
    // OPTIONAL FIELD:
    "replyToId": "MessageId"
} 
```
Outcome:  
```json
{  
    "messageId": "MessageId UUIDv7"  
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
    "createdAt": "Time when it was created",
    "members": [
        {
            "userId": "UUID", 
            "name": "Name",
            "joinedAt": "Time"
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
    "message": "Message text P1",
    "senderKey": "key",
    "receiverKey": "key",
    "nonce": "nonce"
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
    "status": 1,
    "lastSeen": 1785881345
}
```
Statuses:  
- Offline = 1  
- Online = 2  
- Away = 3  
- Typing = 4  

LastSeen - Unix time  
If LastSeen = 0 then profile not exists or was online a very long time ago  