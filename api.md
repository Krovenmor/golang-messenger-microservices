# **Auth API:**  
## **"POST /api/auth/send-code"** - send code to user email (needs for register)
### Request body:
```json
{
    "email": "user@email.com"
}  
```
### Response body:
```json
{
    "expiresIn": 300,   // In seconds
	"retryAfter": 60    // In seconds
}  
```

## **"POST /api/auth/register"** - register a new user  
### Request body:
```json
{
    "login": "YourLogin",
    "password": "YourPassword",
    "email": "user@email.com",
    "code": "verification code from /send-code"
}  
```
- On success code: 201

## **"POST /api/auth/login"** - LogIn, get JWT tokens  
### Request body:
```json
{
    "login": "YourLogin",
    "password": "YourPassword"
}  
```
### Response body  
```json
{
    "accessToken": "YourAccessToken",
    "refreshToken": "YourRefreshToken"
}  
```

## **"POST /api/auth/update"** - Update your tokens  
### Request body:
```json
{
    "refreshToken": "YourRefreshToken"
}  
```
### Response body  
```json
{
    "accessToken": "YourNewAccessToken",
    "refreshToken": "YourNewRefreshToken"
}  
```

## **"GET /api/auth/account"** - Get info about your account 
- You must provide access token in Authorization
### Response body:
```json
{
    "login": "YourLogin",
    "email": "user@email.com"
}
```



# **Profile API:**
**All endpoints works only with JWT token from Auth API**  
**Usage: Authorization -> Bearer Token -> Access token from AuthApi**  

## **"POST /api/profile"** - creates a new profile  
### Request body:  
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

## **"GET /api/profile"** - get yours profile
### Response body (PRIVATE RESPONSE)  
```json
{  
    "userId": "UUID of profile",
    "name": "Name",
    "userName": "User Name",  
    "pubKey": "Public key",  
    "encryptedPrvKey": "Encrypted Private Key",  
    "kdfSalt": "KDF Salt",
    "keyNonce": "Your nonce",
    "createdAt": "time, example: '2026-07-28T19:56:51.855208Z'",

    // Additional Info, {} if no additional info
    "additional": {
        // [] if no avatars
        "avatars": ["UUID from Media Api", ...]
    }
}  
```

## **"GET /api/profile/{target}"** - get other profile  
{target} - UserId (UUID) or UserName  
### Response body (PUBLIC RESPONSE)  
```json
{  
    "userId": "UUID of profile",  
    "name": "Name",  
    "userName": "User Name",  
    "pubKey": "Public key",
    "createdAt": "time, example: '2026-07-28T19:56:51.855208Z'",

    "additional": {
        "avatars": ["UUID from Media Api", ...]
    }
}  
```

## **"POST /api/profile/batch"** - get other profiles  
- limits: from 1 to 50
### Request body:  
```json
{  
    "profiles": ["UUID_1", "UUID_2", ...]
}  
```
### Response body  
```json
[  
    PUBLIC_RESPONSE_1, PUBLIC_RESPONSE_2, ... 
]  
```

## **"POST /api/profile/avatar/{uuid}"** - add new avatar to your profile
- uuid - Your photoId from Media API

## **"DELETE /api/profile/avatar/{uuid}"** - del avatar from your profile
- uuid - Your photoId from Media API



# **Message API:**
**All endpoints works only with JWT token from Auth API**  
**Usage: Authorization -> Bearer Token -> Access token from AuthApi**  

### MessageResponseBody:  
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
    "replyToId": "MessageId UUIDv7",
    // can be empty (not null)
    "reactions": [
        {
            "emoji": "🌈",
            "users": [ // Who reacted
                "98d199b7-0513-4418-8acd-ccb3d4d67575"
            ]
        }
    ]
}
```

## **"GET /api/msg/chats"** - get all your chats  
### Response body  
```json
[  
    "UUID1", "UUID2", ...  
]  
```

## **"GET /api/msg/chats/full"** - get all your chats with extended info  
### Response body  
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
        "lastMessage": MessageResponseBody // (can be null), without reactions !
    }, ...
]  
```  

## **"POST /api/msg/chat"** - create a new chat  
### Request body:  
```json
{
    "userId": "User UUID to create with"  
}  
```
### Response body  
```json
{  
    "chatId": "Chat UUID"  
}  
```
- If chat already exists, you will recv ResponseBody and StatusOk 200, if chat was created: ResponseBody and StatusCreated 201

## **"POST /api/msg/chat/{ChatId}"** - post a message to a chat  
### Request body:  
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
### Response body  
```json
{  
    "messageId": "MessageId UUIDv7"  
}
```

## **"GET /api/msg/chat/{ChatId}"** - get chat history  
Query params:  
'from' - messageId from (UUIDv7) (if you don't provide 'from', you'll get last 'q' messages)
'q' - quantity (int) (if you don't provide 'q', q will be default = 10)
### Response body  
```json
[  
    MessageResponseBody, ...
]  
```

## **"GET /api/msg/chat/{ChatId}/info"** - get chat info  
### Response body  
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

## **"GET /api/msg/chat/{ChatId}/message/{MessageId}"** - get message info  
### Response body  
```json
MessageResponseBody 
```

## **"PUT /api/msg/chat/{ChatId}/message/{MessageId}"** - change message text  
### Request body:  
```json
{   
    "message": "Message text P1",
    "senderKey": "key",
    "receiverKey": "key",
    "nonce": "nonce"
} 
```

## **"DELETE /api/msg/chat/{ChatId}/message/{MessageId}"** - delete message from chat  

## **"GET /api/msg/reactions"** - get all supported reactions
Response:
```json
[
    "🔥", // Emojis
    "💯",
    ...
]
```

## **"POST /api/msg/chat/{ChatId}/message/{MessageId}/reaction"** - post a reaction on message
Request:
```json
{
    "emoji": "emoji only from /api/msg/reactions"
}
```  

## **"DELETE /api/msg/chat/{ChatId}/message/{MessageId}/reaction/{emoji}"** - del a reaction on message
- {emoji} - real emoji, example: /reaction/🔥



# **WebSocket API:**  
**Endpoint works only with JWT token from Auth API**  
## **"WS: /api/ws?token={YourAccessToken}"** - start websocket connection  
By connecting to socket, you will recieve events:  
```json
{  
    "type": "EventType",  
    "payload": EventPayload  
}  
```

## **Events:**  
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
- NewReactionEvent:
```json
{
    "type": "newReaction",
    "payload": {
        "chatId": "UUIDv4",
        "msgId": "UUIDv7",
        "userId": "UUIDv4",
        "emoji": "🌈" // emoji from /reactions endpoint
    }
} 
```
- DelReactionEvent:
```json
{
    "type": "delReaction",
    "payload": {
        "chatId": "UUIDv4",
        "msgId": "UUIDv7",
        "userId": "UUIDv4",
        "emoji": "🌈" // emoji from /reactions endpoint
    }
} 
```

## **Statuses:**  
- You can send your status via ws connection:  
### Request body: (Example)  
```json
{
    "req": 1,
    "payload": {
        "newStatus": 3
    }
}
```
### Response body (Example)  
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

## You can also make request to track someone statuses (primary use: 1v1 chat):  
### Request body: (Example)  
```json
{
    "req": 2,
    "payload": {
        "userId": "{{OtherUserId}}"
    }
}
```
### Response body (Example)  
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
### Response body (Example)
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



# **Media API:**
## - "GET /api/media/profile" - get info
### Response body (Example)  
```json
{
    // In bytes
    "maxSpace": 1234,
    // In bytes
	"spaceFilled": 123,
    "savedFiles": 2
}
```

## - "GET /api/media/profile/files" - get saved files info
- Query params:  
'from' - mediaId from (UUIDv7) (if you don't provide 'from', you'll get last 'q' files)
'q' - quantity (int, default = 10, min=1, max=100)
### Response body (Example)  
```json
[
    {
        "mediaId": "UUIDv7",
        "type": "photo",
        "subType": "avatar",
        // In bytes
        "size": 123,
        "isPublic": true/false,
        // (RFC 3339 format)
        "addedAt": "2026-07-28T00:21:09.904497Z"
    }
]
```
- All possible types:
    - "photo", sub types:
        - "avatar"


## - "POST /api/media/public/avatar" - upload new avatar photo
### Request body: (Example)  
`multipart/form-data` by key 'image'
### Response body (Example)  
```json
{
    "photoId": "UUIDv7"
}
```
### Restrictions:
- Max image size: 600x600 px
- Supported formats: png, jpeg

## - "DELETE /api/media/public/avatar/{photoId}" - delete your avatar photo

# **Static Vault:**
## - "GET /static/avatars/{photoId}" - get avatar (it's public and not encrypted)
### Response will be image/jpeg