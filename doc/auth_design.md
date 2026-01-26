# Authentication & Connection Design

Since the standard browser WebSocket API does not allow setting custom HTTP headers (like `Authorization`) during the initial handshake, we will use a **session-token authentication strategy** via URL query parameters.

## Authentication Flow

1.  **Login (HTTP):** The client authenticates via a standard HTTP POST request.
2.  **Token Issuance:** The server verifies credentials and creates a session record in the database with an opaque token and expiry.
3.  **Connection (WS):** The client opens a WebSocket connection, passing the token in the query string.
4.  **Verification:** The server looks up the session token in the database before upgrading the connection.

---

## 1. Login Endpoint

**URL:** `POST /api/login`
**Content-Type:** `application/json`

**Request Body:**
```json
{
  "username": "user123",
  "password": "secret_password" 
}
```
*(Note: For this prototype, we may just require a username if strictly anonymous/ephemeral, but a password field is good practice for structure.)*

**Response (200 OK):**
```json
{
  "token": "opaque_session_token",
  "expires_in": 3600
}
```

**Response (401 Unauthorized):** Invalid credentials.

---

## 2. WebSocket Connection

**URL:** `ws://<server_host>:<port>/ws?token=<session_token>`

The client appends the received token to the query string.

### Server-Side Handshake Logic

1.  **Intercept:** The Go HTTP handler for `/ws` receives the request.
2.  **Extract:** Parse `token` from the query parameters.
3.  **Validate:** 
    *   Look up the token in the `sessions` table.
    *   Check expiration (`expires_at`).
4.  **Upgrade:** 
    *   **If Valid:** Call `websocket.Upgrader.Upgrade` to establish the socket. Load user info via the session's `user_id` and attach it to the internal Client struct.
    *   **If Invalid:** Return HTTP 401 Unauthorized immediately; do not upgrade.
