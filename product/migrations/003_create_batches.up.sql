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