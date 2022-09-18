# Reserva public API

All routes require bearer token authentication unless otherwise noted.

## Banks
---

### `/v1/banks`
- `GET`
  - Gets the requesting bank's info.
  ### ***Request***
  `GET` with no query parameters.
  
  ### ***Response***
  ```
  {
    "id": <number>,
    "external_id": <number>,
    "email": <string>,
    "source_ips": [
      <string...ip address>,
      <string...ip address>...
    ],
    "name": <string>,
    "balance_in_cents": <number>,
    "frozen": <boolean>
  }
  ```

### `/v1/banks/activate`
- `PATCH`
  - Activates the requesting bank.
  - Bearer token not required. Client authenticates with a token that is delivered by email.
  ### ***Request***
  ```
  {
    "token": <string>
  }
  ```
  ### ***Response***
  ```
  {
    "id": <number>,
    "external_id": <number>,
    "email": <string>,
    "source_ips": [
      <string...ip address>,
      <string...ip address>...
    ],
    "name": <string>,
    "balance_in_cents": <number>,
    "frozen": <boolean>
  }
  ```

### `/v1/banks/update-password`
- `PATCH`
  - Changes the requesting bank's password.
  - Bearer token not required. Client authenticates with a token that is delivered by email.
  ### ***Request***
  ```
  {
    "password": <string>,
    "token": <string>
  }
  ```
  ### ***Response***
  ```
  {
    "id": <number>,
    "external_id": <number>,
    "email": <string>,
    "source_ips": [
      <string...ip address>,
      <string...ip address>...
    ],
    "name": <string>,
    "balance_in_cents": <number>,
    "frozen": <boolean>
  }
  ```

## Accounts

### `/v1/accounts`
- `GET`
  - Gets an accounts info
  ### ***Request***
  `GET` with `?id=<number>` query parameter.
  ### ***Response***
  ```
  {
    "id": <number>,
    "creating_bank": <number>,
    "frozen": <boolean>,
    "balance_in_cents": <number>,
    "depositors": [
      <number...despositor's id>,
      <number...despositor's id>...
    ]
  }
  ```
- `POST`
  - Creates an account by specifying the depositors.
  ### ***Request***
  ```
  [
    <number...despositor's id>,
    <number...despositor's id>...
  ]
  ```
  ### ***Response***
  ```
  {
    "id": <number>,
    "creating_bank": <number>,
    "frozen": <boolean>,
    "balance_in_cents": <number>,
    "depositors": [
      <number...despositor's id>,
      <number...despositor's id>...
    ]
  }
  ```
- `PUT`
  - Changes account depositors.
  - Must submit *all* depositors that should be attached to the account moving forward, including depositors already attached to account.
  ### ***Request***
  ```
  {
    "id": <number>,
    "depositors: [
      <number...despositor's id>,
      <number...despositor's id>...
    ]
  }
  ```
  ### ***Response***
  ```
  {
    "id": <number>,
    "creating_bank": <number>,
    "frozen": <boolean>,
    "balance_in_cents": <number>,
    "depositors": [
      <number...despositor's id>,
      <number...despositor's id>...
    ]
  }
  ```

### `/v1/accounts/freeze`
- `PATCH`
  - Freezes or unfreezes an account.
  ### ***Request***
  ```
  {
    "id": <number>,
    "frozen": <boolean>
  }
  ```
  ### ***Response***
  ```
  {
    "id": <number>,
    "creating_bank": <number>,
    "frozen": <boolean>,
    "balance_in_cents": <number>,
    "depositors": [
      <number...despositor's id>,
      <number...despositor's id>...
    ]
  }
  ```

## Transfers

### `/v1/transfers`
- `POST`
  - Create a new transfer.
  ### ***Request***
  ```
  {
    "source_account": <number>,
    "target_account": <number>,
    "amount_in_cents": <number>
  }
  ```
  ### ***Response***
  ```
  {
    "id": <number>,
    "created_at": <string...RFC 3339>
    "source_account": <number>,
    "target_account": <number>,
    "amount_in_cents": <number>
  }
  ```
- `GET`
  - Allows you to search for a transaction with various filtering critera.
  ### ***Request***
  `GET` with either `id` or one of `source_account` or `target_account`. Also can filter with `created_at`.
  ### ***Response***
  ```
  [
    {
      "id": <number>,
      "created_at": <string...RFC 3339>
      "source_account": <number>,
      "target_account": <number>,
      "amount_in_cents": <number>
    },
    {
      "id": <number>,
      "created_at": <string...RFC 3339>
      "source_account": <number>,
      "target_account": <number>,
      "amount_in_cents": <number>
    }...
  ]
  ```

## Tokens

### `/v1/tokens/authentication`
- `POST`
  - Create an authentication token
  ### ***Request***
  ```
  {
    "email": <string>,
    "password": <string>
  }
  ```
  ### ***Response***
  ```
  {
    "authentication_token": <string>
  }
  ```

### `/v1/tokens/activate`
- `POST`
  - Create an activation token
  ### ***Request***
  ```
  {
    "email": <string>
  }
  ```
  ### ***Response***
  ```
  {
    "message": <string>
  }
  ```

### `/v1/tokens/reset-password`
- `POST`
  - Create a password reset token
  ### ***Request***
  ```
  {
    "email": <string>
  }
  ```
  ### ***Response***
  ```
  {
    "message": <string>
  }
  ```

## Utility

### `/v1/healthcheck`
- `GET`
  - A healthcheck
  ### ***Request***
  `GET` with no query parameters.
  ### ***Response***
  ```
  {
		"status": <string>,
		"system_info": {
			"environment": <string>,
			"version": <string
		}
	}
  ```