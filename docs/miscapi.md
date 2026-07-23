# Misc API

I don't care enough about these, so I shoved them all here. They aren't all too useful.

## Healthz Check

Returns `OK` in plain text.

- **HTTP Method:** `GET`
- **Path:** `/api/healthz`

### Example Request

```bash
curl -X GET "https://example.com/api/healthz"
```

### Success Response

**Status:** `200 OK`

## Request Amount

Returns html that tells you how many hits the file server handler has.

- **HTTP Method:** `GET`
- **Path:** `/admin/metrics`

### Example Request

```bash
curl -X GET "https://example.com/admin/metrics"
```

### Success Response

**Status:** `200 OK`

```html
<html>
<body>
<h1>Welcome, Chirpy Admin</h1>
<p>Chirpy has been visited 5 times!</p>
</body>
</html>
```

## Reset Database

Deletes all records from the database. Also resets the hits on the file server.

- **HTTP Method:** `POST`
- **Path:** `/admin/reset`

### Example Request

```bash
curl -X POST "https://example.com/admin/reset"
```

### Success Response

**Status:** `200 OK`

### Error Responses

- **`403 Forbidden`**: Returned when `PLATFORM` in the .env file doesn't equal `dev`.
  ```json
  {
    "error": "no access to endpoint"
  }
  ```
