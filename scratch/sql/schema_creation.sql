drop table if exists tokens;

drop table if exists scope_types;

drop table if exists transfers;

drop table if exists cards;

drop table if exists accounts;

drop table if exists depositors;

drop table if exists bank_ip;

drop table if exists banks;

create table banks (
  id bigint primary key,
  created_at timestamp not null default now(),
  name text not null unique,
  email text not null unique,
  password_hash bytea not null,
  activated bool not null default false,
  version int not null default 0,
  external_id text not null unique,
  balance_in_cents bigint not null default 0,
  frozen boolean not null default false
);

create table bank_ip (
  bank_id bigint not null references banks on delete cascade,
  ip inet not null
);

create table depositors (
  id bigint primary key,
  external_id text not null unique,
  name text not null,
  address text not null unique,
  date_of_birth date,
  version int not null default 0
);

create table accounts (
  id bigint primary key,
  bank_id bigint references banks,
  frozen boolean not null default false,
  balance_in_cents bigint not null default 0,
  version int not null default 0
);

create table cards (
  id bigint primary key,
  account_id bigint not null references accounts,
  depositor bigint not null references depositors,
  public_key bytea not null,
  pin_hash bytea not null,
  expiry timestamp not null,
  version int not null default 0
);

create table transfers (
  id uuid DEFAULT gen_random_uuid() primary key,
  card_id bigint not null references cards,
  target_account_id bigint not null references accounts,
  created_at timestamp not null default now(),
  completed_at timestamp,
  amount_in_cents bigint not null
);

create table scope_types (
  id bigint generated always as identity primary key,
  name text
);

create table tokens (
  secret_hash bytea primary key,
  bank_id bigint not null references banks on delete cascade,
  expiry timestamp not null,
  scope bigint not null references scope_types
);