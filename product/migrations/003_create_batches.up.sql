-- Create batches table
-- id, product_id, batch_number, quantity, purchase_price, selling_price, expiration_date, supplier_id, purchase_id, created_at, updated_at
CREATE TABLE batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id),
    batch_number VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    purchase_price DECIMAL(10,2) NOT NULL CHECK (purchase_price >= 0),
    selling_price DECIMAL(10,2) NOT NULL CHECK (selling_price >= 0),
    expiration_date DATE NOT NULL,
    supplier_id UUID,
    purchase_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
-- idx_batches_product_id on product_id
-- idx_batches_expiration_date on expiration_date
-- idx_batches_batch_number on batch_number
-- idx_batches_supplier_id on supplier_id
-- idx_batches_purchase_id on purchase_id
CREATE INDEX idx_batches_product_id ON batches(product_id);
CREATE INDEX idx_batches_expiration_date ON batches(expiration_date);
CREATE INDEX idx_batches_batch_number ON batches(batch_number);
CREATE INDEX idx_batches_supplier_id ON batches(supplier_id);
CREATE INDEX idx_batches_purchase_id ON batches(purchase_id);

-- Insert dummy data (5 batches)
INSERT INTO batches (product_id, batch_number, quantity, purchase_price, selling_price, expiration_date) VALUES
((SELECT id FROM products WHERE name = 'Paracetamol 500mg' LIMIT 1), 'BATCH-001', 500, 4000.00, 5000.00, CURRENT_DATE + INTERVAL '2 years'),
((SELECT id FROM products WHERE name = 'Vitamin C 1000mg' LIMIT 1), 'BATCH-002', 200, 12000.00, 15000.00, CURRENT_DATE + INTERVAL '1 year'),
((SELECT id FROM products WHERE name = 'Bandage Roll 5cm' LIMIT 1), 'BATCH-003', 1000, 6000.00, 8000.00, CURRENT_DATE + INTERVAL '3 years'),
((SELECT id FROM products WHERE name = 'Face Mask 3-Ply' LIMIT 1), 'BATCH-004', 2000, 2000.00, 3000.00, CURRENT_DATE + INTERVAL '2 years'),
((SELECT id FROM products WHERE name = 'Baby Shampoo 200ml' LIMIT 1), 'BATCH-005', 150, 20000.00, 25000.00, CURRENT_DATE + INTERVAL '1 year');