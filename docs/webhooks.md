# Webhook API

Webhook endpoints for services.

## Authentication

All requests require an apikey in the HTTP `Authorization` header. If you want to use the single webhook, use the `POLKA_KEY` in the .env file as the apikey.

```http
Authorization: Bearer <APIKEY>
```

## Polka

Exucutes any given Polka event. Only the `user.upgraded` event does anything.

- **HTTP Method:** `POST`
- **Path:** `/api/polka/webhooks`
- **Content-Type:** `application/json`

### Request Body

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `event` | string | **Yes** | The type of event. |
| `data` | object | **Yes** | Payload container for event details. |
| `data.user_id` | uuid | **Yes** | The UUID of the user affected by the event. |

### Example Request

```bash
curl -X POST "https://example.com/api/polka/webhooks" \
  -H "Content-Type: application/json" \
  -d '{
    "event": "user.upgraded",
    "data": {
        "user_id": "user-uuid"
    }
  }'
```

### Success Response

**Status:** `204 No Content`

Note: Status will also be 204 if event is not `user.upgraded`, but will do nothing.

### Error Responses

- **`400 Bad Request`**: Returned when the request body cannot be decoded as JSON.
  ```json
  {
    "error": "error decoding request body"
  }
  ```
- **`401 Unauthorized`**: Returned when the `Authorization` header is missing.
  ```json
  {
    "error": "No Authorization header found"
  }
  ```
- **`401 Unauthorized`**: Returned when the apikey is wrong.
  ```json
  {
    "error": "wrong apikey"
  }
  ```
- **`404 Not Found`**: Returned when user does not exist.
  ```json
  {
    "error": "user not found"
  }
  ```
- **`500 Internal Server Error`**: Returned when a server error occurs.
  ```json
  {
    "error": "An unexpected server error occurred"
  }
  ```
