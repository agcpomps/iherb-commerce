CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  sku TEXT UNIQUE NOT NULL,
  category TEXT,
  margin_percent NUMERIC NOT NULL CHECK (margin_percent >= 0),
  active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  image_path TEXT
);

CREATE TABLE boxes (
  id UUID PRIMARY KEY,
  supplier_name TEXT NOT NULL,
  purchase_date DATE NOT NULL,
  exchange_rate_usd_aoa NUMERIC NOT NULL CHECK (exchange_rate_usd_aoa > 0),
  shipping_usd NUMERIC NOT NULL DEFAULT 0,
  customs_tax_usd NUMERIC NOT NULL DEFAULT 0,
  other_costs_usd NUMERIC NOT NULL DEFAULT 0,
  status TEXT NOT NULL CHECK (status IN ('in_transit','arrived')),
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE batches (
  id UUID PRIMARY KEY,
  product_id UUID REFERENCES products(id),
  box_id UUID REFERENCES boxes(id),
  quantity_received INT NOT NULL CHECK (quantity_received > 0),
  quantity_available INT NOT NULL CHECK (quantity_available >= 0),
  expiration_date DATE NOT NULL,
  landed_cost_usd NUMERIC NOT NULL CHECK (landed_cost_usd > 0),
  landed_cost_aoa NUMERIC NOT NULL CHECK (landed_cost_aoa > 0),
  status TEXT NOT NULL CHECK (status IN ('active','near_expiry','expired')),
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_batches_fifo
  ON batches (product_id, expiration_date);
