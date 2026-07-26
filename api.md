**Auth API:**  
-

- "POST /api/auth/register", h.Register - register a new user  
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