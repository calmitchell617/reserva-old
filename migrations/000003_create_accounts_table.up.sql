create table accounts (
  id bigserial primary key,
  bank_id bigint references banks,
  frozen boolean not null default false,
  balance_in_cents bigint not null default 0,
  version bigint not null default 0
);