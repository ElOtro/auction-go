CREATE EXTENSION IF NOT EXISTS "citext";
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  active bool DEFAULT false,
  role integer DEFAULT 1,
  name text NOT NULL,
  email citext UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) without time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX users_active_index ON users USING btree (active);
CREATE UNIQUE INDEX users_email_index ON users USING btree (email);

comment on column users.active is 'Active/Disabled';
comment on column users.role is 'Role';
comment on column users.name is 'Name';
comment on column users.email is 'Email';
comment on column users.password_hash is 'Hashed Password';