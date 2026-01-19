# Database Schema Design

For this application, we will use a relational database structure. This schema is compatible with **SQLite** (for ease of development/embedded use) or **PostgreSQL** (for production).

## Users Table

The `users` table handles identity and authentication credentials.

**Table Name:** `users`

| Column Name | Data Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` or `TEXT` | **PK**, Not Null | Unique identifier for the user. |
| `username` | `VARCHAR(50)` | **Unique**, Not Null | The display name used for login and chat. |
| `password_hash`| `VARCHAR(255)` | Not Null | The **bcrypt** hash of the user's password. *Never store plain text.* |
| `created_at` | `TIMESTAMP` | Default: `NOW()` | When the account was registered. |
| `last_seen` | `TIMESTAMP` | Nullable | Timestamp of the user's last activity/login. |

### SQL Definition (PostgreSQL Example)

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP WITH TIME ZONE
);

-- Index for fast lookups during login
CREATE INDEX idx_users_username ON users(username);
```

### Go Struct Mapping (GORM)

If using GORM (Go Object Relational Mapper), the model would look like this:

```go
type User struct {
    ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` 
    Username     string    `gorm:"uniqueIndex;not null;size:50"`
    PasswordHash string    `gorm:"not null"`
    CreatedAt    time.Time `gorm:"autoCreateTime"`
    LastSeen     time.Time
}
```
