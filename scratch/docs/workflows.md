# Reserva workflows

## New bank signup
- Customer goes through KYC
- Commercial bank creates a depositor profile
- Customer accesses Reserva through the bank's client applications

## New account creation
- Commercial bank creates a new account and assigns depositors

## Card generation
- Customer inputs a pin
- A new record is created in cards table
- Input by bank
  - account_id bigint fk
  - depositor_id bigint fk
  - hashed_pin bytea
- Generated
  - id bigint pk
  - public_key bytea
- Private key and id are flashed onto the card


## Transaction is requested
- Payment service provider (PSP) hits new transaction route with a post
  - Target account
  - Card ID
- Server validates input, checks balance, and creates a transfer in the DB.
  - Hashes id (uuid) with card's public key, sends to PSP
- PSP decrypts secret by running it through the card
- PSP sends decrypted secret back to server, which completes transaction, and sends 200 response with transaction info