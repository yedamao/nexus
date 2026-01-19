# Authentication & Connection Design

Since the standard browser WebSocket API does not allow setting custom HTTP headers (like `Authorization`) during the initial handshake, we will use a **Ticket/Token-based authentication strategy** via URL query parameters.

## Authentication Flow

1.  **Login (HTTP):** The client authenticates via a standard HTTP POST request.
2.  **Token Issuance:** The server verifies credentials and issues a signed, short-lived **JWT (JSON Web Token)**.
3.  **Connection (WS):** The client opens a WebSocket connection, passing the JWT in the query string.
4.  **Verification:** The server validates the token before upgrading the connection.

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
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600
}
```

**Response (401 Unauthorized):** Invalid credentials.

---

## 2. WebSocket Connection

**URL:** `ws://<server_host>:<port>/ws?token=<jwt_token>`

The client appends the received token to the query string.

### Server-Side Handshake Logic

1.  **Intercept:** The Go HTTP handler for `/ws` receives the request.
2.  **Extract:** Parse `token` from the query parameters.
3.  **Validate:** 
    *   Parse the JWT.
    *   Verify the signature (using the server's secret key).
    *   Check expiration (`exp` claim).
4.  **Upgrade:** 
    *   **If Valid:** Call `websocket.Upgrader.Upgrade` to establish the socket. Extract user info from claims (e.g., `sub` or `username`) and attach it to the internal Client struct.
    *   **If Invalid:** Return HTTP 401 Unauthorized immediately; do not upgrade.

---

## Token Payload (JWT Claims)

The JWT payload will contain the minimal necessary user context:

```json
{
  "sub": "user_uuid_123",   // Subject (User ID)
  "username": "user123",    // Display Name
  "iat": 1698393000,        // Issued At
  "exp": 1698396600         // Expiration (e.g., 1 hour)
}
```
