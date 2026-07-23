# Users API

Use these endpoints to manage user accounts and authentication-related profile data.

Note: These endpoints do not enforce email uniqueness or password complexity rules.

## Authentication

Some requests require an access token in the HTTP `Authorization` header. Use the login endpoint to obtain a token first.

```http
Authorization: Bearer <YOUR_ACCESS_TOKEN>
```

## Create a User

Creates a new user account.

- **HTTP Method:** `POST`
- **Path:** `/api/users`
- **Content-Type:** `application/json`

### Request Body

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `email` | string | **Yes** | A non-empty email address. |
| `password` | string | **Yes** | A non-empty password. |

### Example Request

```bash
curl -X POST "https://example.com/api/users" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane.doe@example.com",
    "password": "password123"
  }'
```

### Success Response

**Status:** `201 Created`

```json
{
  "id": "generated-uuid",
  "created_at": "2026-07-23T04:24:00Z",
  "updated_at": "2026-07-23T04:24:00Z",
  "email": "jane.doe@example.com",
  "is_chirpy_red": false
}
```

### Error Responses

- **`400 Bad Request`**: Returned when the request body cannot be decoded as JSON.
  ```json
  {
    "error": "error decoding request body"
  }
  ```
- **`400 Bad Request`**: Returned when either `email` or `password` is empty.
  ```json
  {
    "error": "please fill out both the email and password fields"
  }
  ```
- **`500 Internal Server Error`**: Returned when a server error occurs.
  ```json
  {
    "error": "An unexpected server error occurred"
  }
  ```

## Login as User

Authenticates a user and returns an access token plus a refresh token. See `refreshapi.md` for how to use the refresh_token.

- **HTTP Method:** `POST`
- **Path:** `/api/login`
- **Content-Type:** `application/json`

### Request Body

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `email` | string | **Yes** | The user's email address. |
| `password` | string | **Yes** | The user's password. |

### Example Request

```bash
curl -X POST "https://example.com/api/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane.doe@example.com",
    "password": "password123"
  }'
```

### Success Response

**Status:** `200 OK`

```json
{
  "id": "user-uuid",
  "created_at": "2026-07-23T04:24:00Z",
  "updated_at": "2026-07-23T04:24:00Z",
  "email": "jane.doe@example.com",
  "is_chirpy_red": false,
  "token": "new-access-token",
  "refresh_token": "new-refresh-token"
}
```

### Error Responses

- **`400 Bad Request`**: Returned when the request body cannot be decoded as JSON.
  ```json
  {
    "error": "error decoding request body"
  }
  ```
- **`401 Unauthorized`**: Returned when the email does not exist or the password is incorrect.
  ```json
  {
    "error": "Incorrect email or password"
  }
  ```
- **`500 Internal Server Error`**: Returned when a server error occurs.
  ```json
  {
    "error": "An unexpected server error occurred"
  }
  ```

## Change User Email and Password

Updates a user's email and password. This endpoint requires a valid access token.

- **HTTP Method:** `PUT`
- **Path:** `/api/users`
- **Content-Type:** `application/json`

### Request Body

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `email` | string | **No** | A non-empty email address. |
| `password` | string | **No** | A non-empty password. |

### Example Request

```bash
curl -X PUT "https://example.com/api/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <YOUR_ACCESS_TOKEN>" \
  -d '{
    "email": "jane.doe@example.com",
    "password": "password123"
  }'
```

### Success Response

**Status:** `200 OK`

```json
{
  "id": "user-uuid",
  "created_at": "2026-07-23T04:24:00Z",
  "updated_at": "2026-07-26T07:35:00Z",
  "email": "jane.doe@example.com",
  "is_chirpy_red": false
}
```

### Error Responses

- **`400 Bad Request`**: Returned when the request body cannot be decoded as JSON.
  ```json
  {
    "error": "error decoding request body"
  }
  ```
- **`400 Bad Request`**: Returned when both `email` or `password` is empty.
  ```json
  {
    "error": "at least the email or password fields must be filled out"
  }
  ```
- **`401 Unauthorized`**: Returned when the `Authorization` header is missing.
  ```json
  {
    "error": "No Authorization header found"
  }
  ```
- **`401 Unauthorized`**: Returned when the access token is invalid.
  ```json
  {
    "error": "JWT token is invalid"
  }
  ```
- **`500 Internal Server Error`**: Returned when a server error occurs.
  ```json
  {
    "error": "An unexpected server error occurred"
  }
  ```