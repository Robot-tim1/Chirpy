# Refresh API

Use these endpoints to get new access tokens, and revoke a refresh token.

## Authentication

All requests require a refresh token in the HTTP `Authorization` header. See `usersapi.md` for how to get a refresh token.

```http
Authorization: Bearer <YOUR_REFRESH_TOKEN>
```

## Refresh Access Token

Creates a new access token to be used. Require a refresh token in the HTTP `Authorization` header. Does not accept request bodies.

- **HTTP Method:** `POST`
- **Path:** `/api/refresh`

### Example Request

```bash
curl -X POST "https://example.com/api/refresh" \
  -H "Authorization: Bearer <YOUR_REFRESH_TOKEN>"
```

### Success Response

**Status:** `200 OK`

```json
{
  "token": "new-access-token"
}
```

### Error Responses

- **`401 Unauthorized`**: Returned when the `Authorization` header is missing.
  ```json
  {
    "error": "No Authorization header found"
  }
  ```
- **`401 Unauthorized`**: Returned when the refresh token does not exist in the database.
  ```json
  {
    "error": "refresh token does not exist"
  }
  ```
- **`401 Unauthorized`**: Returned when the refresh token has been revoked.
  ```json
  {
    "error": "refresh token is invalid"
  }
  ```
- **`500 Internal Server Error`**: Returned when a server error occurs.
  ```json
  {
    "error": "An unexpected server error occurred"
  }
  ```

## Revoke Refresh Token

Revokes a given refresh token. Require a refresh token in the HTTP `Authorization` header. Does not accept request bodies.

- **HTTP Method:** `POST`
- **Path:** `/api/revoke`

### Example Request

```bash
curl -X POST "https://example.com/api/revoke" \
  -H "Authorization: Bearer <YOUR_REFRESH_TOKEN>"
```

### Success Response

**Status:** `204 No Content`

### Error Responses

- **`401 Unauthorized`**: Returned when the `Authorization` header is missing.
  ```json
  {
    "error": "No Authorization header found"
  }
  ```
- **`500 Internal Server Error`**: Returned when a server error occurs.
  ```json
  {
    "error": "An unexpected server error occurred"
  }
  ```