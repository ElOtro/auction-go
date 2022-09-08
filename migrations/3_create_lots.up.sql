CREATE TABLE lots (
  id BIGSERIAL PRIMARY KEY,
  status integer DEFAULT 1,
  title text,
  description text,
  start_price bigint DEFAULT 0,
  end_price bigint DEFAULT 0,
  step_price bigint DEFAULT 0,
  creator_id bigint REFERENCES users (id) ON DELETE SET NULL,
  winner_id bigint REFERENCES users (id) ON DELETE SET NULL,
  start_at timestamp(0) with time zone,
  end_at timestamp(0) with time zone,
  notify bool DEFAULT false,
  destroyed_at timestamp(0) without time zone,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) without time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX lots_status_index ON lots USING btree (status);
CREATE INDEX lots_start_at_index ON lots USING btree (start_at);
CREATE INDEX lots_end_at_index ON lots USING btree (end_at);
CREATE INDEX lots_destroyed_at_index ON lots USING btree (destroyed_at);

comment on column lots.status is 'Status';
comment on column lots.title is 'Name (Number)';
comment on column lots.description is 'Description';
comment on column lots.start_price is 'Start Price';
comment on column lots.end_price is 'End Price';
comment on column lots.step_price is 'Step Price';
comment on column lots.creator_id is 'Creator ID (User)';
comment on column lots.winner_id is 'Winner ID (User)';
comment on column lots.start_at is 'Start Datetime (Auction)';
comment on column lots.end_at is 'End Datetime (Auction)';
comment on column lots.notify is 'Notify On Start';