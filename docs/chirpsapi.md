# Chirps API

Use these endpoints to post, get and delete chirps (like tweets but not twitter) from the database.

## Authentication

Some requests require an access token in the HTTP `Authorization` header. See `usersapi.md` for how to get an access token.

```http
Authorization: Bearer <YOUR_ACCESS_TOKEN>
```

## Create a Chirp

Creates a new chirp post in the database. Requires an access token in `Authorization` header.

- **HTTP Method:** `POST`
- **Path:** `/api/chirps`
- **Content-Type:** `application/json`

### Request Body

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `body` | string | **Yes** | A string of text under 140 characters. |

### Example Request

```bash
curl -X POST "https://example.com/api/chirps" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <YOUR_ACCESS_TOKEN>" \
  -d '{
    "body": "chirping on chirp"
  }'
```

### Success Response

**Status:** `201 Created`

```json
{
  "id": "chirp-uuid",
  "created_at": "2026-07-23T04:24:00Z",
  "updated_at": "2026-07-23T04:24:00Z",
  "body": "body-from-chirp",
  "user_id": "author-uuid"
}
```

### Error Responses

- **`400 Bad Request`**: Returned when the request body cannot be decoded as JSON.
  ```json
  {
    "error": "error decoding request body"
  }
  ```
- **`400 Bad Request`**: Returned when the request body text is more than or equal to 140.
  ```json
  {
    "error": "Chirp is too long"
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

## Getting Chirps

Get chirps from the database.

- **HTTP Method:** `GET`
- **Path:** `/api/chirps`

### Parameters

#### Path Parameters

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `chirpID` | uuid | **No** | A uuid of a chirp. Gets the chirp if it exists. |

#### Query Parameters

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `author_id` | uuid | **No** | A uuid of a user. Shows all posts from given user. |
| `sort` | string | **No** | sort either in asc or desc order. Default is asc. |

### Example Request

```bash
curl -X GET "https://example.com/api/chirps"
```

### Success Response

**Status:** `200 OK`

```json
{
  "id": "chirp-uuid",
  "created_at": "2026-07-23T04:24:00Z",
  "updated_at": "2026-07-23T04:24:00Z",
  "body": "body-from-chirp",
  "user_id": "author-uuid"
}
```

### Error Responses

- **`400 Bad Request`**: Returned when `chirpID` path parameter is not a valid uuid.
  ```json
  {
    "error": "error parsing id into uuid"
  }
  ```
- **`400 Bad Request`**: Returned when `sort` query parameter is not asc or desc.
  ```json
  {
    "error": "no sort query of that kind"
  }
  ```
- **`404 Not Found`**: Returned when `author_id` has no chirps.
  ```json
  {
    "error": "author has no posts"
  }
  ```
- **`404 Not Found`**: Returned when `chirpID` does not exist.
  ```json
  {
    "error": "Chirp could not be found"
  }
  ```
- **`500 Internal Server Error`**: Returned when a server error occurs.
  ```json
  {
    "error": "An unexpected server error occurred"
  }
  ```

## Delete Chirp

Delete a chirp from the database. Requires an access token in `Authorization` header.

- **HTTP Method:** `DELETE`
- **Path:** `/api/chirps/{chirpID}`

### Parameters

#### Path Parameters

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `chirpID` | uuid | **Yes** | A uuid of a chirp. Deletes the chirp if it exists. |

### Example Request

```bash
curl -X DELETE "https://example.com/api/chirps/{chirpID}"
```

### Success Response

**Status:** `204 No Content`

### Error Responses

- **`400 Bad Request`**: Returned when `author_id` query parameter is not a valid uuid.
  ```json
  {
    "error": "error parsing id into uuid"
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
- **`403 Forbidden`**: Returned when `chirpID`'s `author_id` does not match user_id from access token.
  ```json
  {
    "error": "Can not delete chirp post from other user"
  }
  ```
- **`404 Not Found`**: Returned when `author_id` has no chirps.
  ```json
  {
    "error": "author has no posts"
  }
  ```
- **`404 Not Found`**: Returned when `chirpID` does not exist.
  ```json
  {
    "error": "Chirp could not be found"
  }
  ```
- **`500 Internal Server Error`**: Returned when a server error occurs.
  ```json
  {
    "error": "An unexpected server error occurred"
  }
  ```