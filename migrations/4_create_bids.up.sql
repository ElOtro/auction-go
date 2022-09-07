CREATE TABLE bids (
  id BIGSERIAL PRIMARY KEY,
  amount bigint DEFAULT 0,
  lot_id bigint REFERENCES lots (id) ON DELETE CASCADE,
  bidder_id bigint REFERENCES users (id) ON DELETE SET NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) without time zone NOT NULL DEFAULT NOW()
);

comment on column bids.amount is 'Amount Of Bid';
comment on column bids.lot_id is 'Lot ID';
comment on column bids.bidder_id is 'Bidder ID (User)';
