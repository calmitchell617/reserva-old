create table banks (
  id bigserial primary key,
  created_at timestamp(0) not null default now(),
  name text not null unique,
  email citext not null unique,
  password_hash bytea not null,
  activated bool not null default false,
  version int not null default 1,
  external_id text not null unique,
  balance_in_cents bigint not null default 0,
  frozen boolean not null default false
);