BEGIN;

CREATE TABLE orders (
  id TEXT PRIMARY KEY,
  total_amount_aoa NUMERIC NOT NULL CHECK (total_amount_aoa >= 0),
  payment_method TEXT NOT NULL
    CHECK (payment_method IN ('transfer_site','transfer_whatsapp')),
  status TEXT NOT NULL
    CHECK (status IN ('pending_payment','paid','cancelled')),
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE order_items (
  id UUID PRIMARY KEY,
  order_id TEXT NOT NULL
    REFERENCES orders(id) ON DELETE CASCADE,
  product_id UUID NOT NULL,
  product_name TEXT NOT NULL,
  batch_id UUID NOT NULL,
  quantity INT NOT NULL CHECK (quantity > 0),
  unit_price_aoa NUMERIC NOT NULL CHECK (unit_price_aoa >= 0),
  total_price_aoa NUMERIC NOT NULL CHECK (total_price_aoa >= 0)
);

CREATE INDEX idx_order_items_order_id
  ON order_items(order_id);

CREATE INDEX idx_order_items_product_id
  ON order_items(product_id);

COMMIT;
