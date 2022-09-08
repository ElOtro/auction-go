CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE accounts (
  id uuid DEFAULT uuid_generate_v4(),
  ammount bigint DEFAULT 0,
  user_id bigint REFERENCES users (id) ON DELETE CASCADE,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) without time zone NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

comment on column accounts.id is 'Account Number';
comment on column accounts.ammount is 'Ammount';
comment on column accounts.user_id is 'User ID';