CREATE TABLE bids (
  id BIGSERIAL PRIMARY KEY,
  price biginteger DEFAULT 0,
  lot_id bigint REFERENCES lots (id) ON DELETE CASCADE,
  bidder_id bigint REFERENCES users (id) ON DELETE SET NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) without time zone NOT NULL DEFAULT NOW()
);
CREATE INDEX companies_destroyed_at_index ON companies USING btree (destroyed_at);