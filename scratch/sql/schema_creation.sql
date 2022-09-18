drop table if exists tokens;

drop table if exists transfers;

drop table if exists cards;

drop table if exists accounts;

drop table if exists banks;

create table banks (
  id bigserial primary key,
  name text not null unique,
  email citext not null unique,
  password_hash bytea not null,
  balance_in_cents bigint not null default 0,
  activated bool not null default false,
  frozen boolean not null default false,
  ips inet[] not null default array[]::inet[],
  version bigint not null default 0
);

create table accounts (
  id bigint primary key,
  bank_id bigint references banks,
  frozen boolean not null default false,
  balance_in_cents bigint not null default 0,
  version bigint not null default 0
);

create table cards (
  id bigint primary key,
  account_id bigint not null references accounts,
  private_key bytea not null,
  pin_hash bytea not null,
  expiry timestamp not null,
  version bigint not null default 0
);

create table transfers (
  id uuid primary key,
  card_id bigint not null references cards,
  target_account_id bigint not null references accounts,
  created_at timestamp not null default now(),
  amount_in_cents bigint not null
);

create table tokens (
  secret_hash bytea primary key,
  bank_id bigint not null references banks on delete cascade,
  expiry timestamp not null,
  scope text not null
);