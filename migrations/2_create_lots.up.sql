CREATE TABLE lots (
  id BIGSERIAL PRIMARY KEY,
  status integer DEFAULT 1,
  name text,
  description text,
  start_price biginteger DEFAULT 0,
  end_price biginteger DEFAULT 0,
  creator_id bigint REFERENCES users (id) ON DELETE SET NULL,
  winner_id bigint REFERENCES users (id) ON DELETE SET NULL,
  start_at timestamp(0) with time zone,
  end_at timestamp(0) with time zone,
  destroyed_at timestamp(0) without time zone,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) without time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX lots_status_index ON lots USING btree (status);
CREATE INDEX lots_start_at_index ON lots USING btree (start_at);
CREATE INDEX lots_end_at_index ON lots USING btree (end_at);
CREATE INDEX lots_destroyed_at_index ON lots USING btree (destroyed_at);