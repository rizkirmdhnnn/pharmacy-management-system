-- Create purchase_items table
-- id, purchase_id, product_id, quantity, total_price, created_at
CREATE TABLE purchase_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    purchase_id UUID NOT NULL REFERENCES purchases(id),
    product_id UUID NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    total_price DECIMAL(10,2) NOT NULL CHECK (total_price >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_purchase_items_purchase_id ON purchase_items(purchase_id);
CREATE INDEX idx_purchase_items_product_id ON purchase_items(product_id);

-- Insert dummy data (5 purchase items - one per purchase)
-- Note: product_id uses placeholder UUIDs since products are in a different service database
INSERT INTO purchase_items (purchase_id, product_id, quantity, total_price) VALUES
((SELECT id FROM purchases WHERE purchase_number = 'PO-001' LIMIT 1), '11111111-1111-1111-1111-111111111111'::uuid, 100, 1500000.00),
((SELECT id FROM purchases WHERE purchase_number = 'PO-002' LIMIT 1), '22222222-2222-2222-2222-222222222222'::uuid, 200, 2400000.00),
((SELECT id FROM purchases WHERE purchase_number = 'PO-003' LIMIT 1), '33333333-3333-3333-3333-333333333333'::uuid, 200, 800000.00),
((SELECT id FROM purchases WHERE purchase_number = 'PO-004' LIMIT 1), '44444444-4444-4444-4444-444444444444'::uuid, 100, 3200000.00),
((SELECT id FROM purchases WHERE purchase_number = 'PO-005' LIMIT 1), '55555555-5555-5555-5555-555555555555'::uuid, 120, 1800000.00);

