# Mini Documentation

## Base URL

`http://<server_address>/api/v1`

## Endpoints

| Method | Endpoint         | Description                                   |
| ------ | ---------------- | --------------------------------------------- |
| POST   | `/account`       | Creates a new account                         |
| GET    | `/account`       | Retrieves account information                 |
| GET    | `/health`        | Health check for the server                   |
| GET    | `/ping`          | Simple ping to check if the server is running |
| POST   | `/token`         | Generates access and refresh tokens           |
| POST   | `/refresh-token` | Refreshes tokens using a valid refresh token  |

## Clarification

### `POST /account`

Request Body

```
{
  "email": "user@example.com",
  "password": "user-password",
  "name": "User Name"
}
```

### `GET /account`

Headers: `Authorization: Bearer <access_token>`

Response body

```
{
  "id": "user-id",
  "email": "user@example.com",
  "name": "User Name",
  "created_at": "timestamp"
}
```

### `GET /health`

Response body

```
{
  "status": "healthy"
}
```

### `GET /ping`

Response body:

```
{
  "message": "pong"
}
```

### `POST /token`

Params: `guid`

Response body:

```
{
  "access_token": "new-access-token",
  "refresh_token": "new-refresh-token"
}

```

### `POST /refresh-token`

Request body:

```
{
  "refresh_token": "valid-refresh-token"
}
```

Response body:

```
{
  "access_token": "new-access-token",
  "refresh_token": "new-refresh-token"
}
```
