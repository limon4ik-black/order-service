CREATE TABLE IF NOT EXISTS general (
  order_uid TEXT PRIMARY KEY,
  track_number TEXT UNIQUE NOT NULL,
  entry TEXT NOT NULL,
  locale TEXT NOT NULL,
  internal_signature TEXT,
  customer_id TEXT NOT NULL,
  delivery_service TEXT NOT NULL,
  shardkey TEXT NOT NULL,
  sm_id integer NOT NULL,
  date_created TIMESTAMP NOT NULL,
  oof_shard TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS delivery (
  id SERIAL PRIMARY KEY,
  order_uid TEXT NOT NULL REFERENCES general(order_uid) ON DELETE CASCADE,
  name  TEXT NOT NULL,
  phone TEXT NOT NULL,
  zip TEXT NOT NULL,
  city TEXT NOT NULL,
  address TEXT NOT NULL,
  region TEXT NOT NULL,
  email TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS payment (
  id SERIAL PRIMARY KEY,
  order_uid TEXT NOT NULL REFERENCES general(order_uid) ON DELETE CASCADE,
  transaction TEXT NOT NULL,
  request_id TEXT,
  currency TEXT NOT NULL,
  provider TEXT NOT NULL,
  amount INT NOT NULL,
  payment_dt BIGINT NOT NULL,
  bank TEXT NOT NULL,
  delivery_cost INT NOT NULL,
  goods_total INT NOT NULL,
  custom_fee INT NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
  id SERIAL PRIMARY KEY,
  order_uid TEXT NOT NULL REFERENCES general(order_uid) ON DELETE CASCADE,
  chrt_id BIGINT NOT NULL,
  track_number TEXT NOT NULL,
  price INT NOT NULL,
  rid TEXT NOT NULL,
  name TEXT NOT NULL,
  sale INT NOT NULL,
  size TEXT NOT NULL,
  total_price INT NOT NULL,
  nm_id BIGINT NOT NULL,
  brand TEXT NOT NULL,
  status INT NOT NULL
);


CREATE OR REPLACE FUNCTION insert_general(
    p_order_uid TEXT,
    p_track_number TEXT,
    p_entry TEXT,
    p_locale TEXT,
    p_internal_signature TEXT,
    p_customer_id TEXT,
    p_delivery_service TEXT,
    p_shardkey TEXT,
    p_sm_id INTEGER,
    p_date_created TIMESTAMP,
    p_oof_shard TEXT
)
RETURNS VOID AS $$
BEGIN
    INSERT INTO general (
        order_uid,
        track_number,
        entry,
        locale,
        internal_signature,
        customer_id,
        delivery_service,
        shardkey,
        sm_id,
        date_created,
        oof_shard
    )
    VALUES (
        p_order_uid,
        p_track_number,
        p_entry,
        p_locale,
        p_internal_signature,
        p_customer_id,
        p_delivery_service,
        p_shardkey,
        p_sm_id,
        p_date_created,
        p_oof_shard
    );
EXCEPTION
    WHEN unique_violation THEN
        RAISE EXCEPTION 'Запись с order_uid "%" или track_number "%" уже существует.',
            p_order_uid, p_track_number;
END;
$$ LANGUAGE plpgsql;
