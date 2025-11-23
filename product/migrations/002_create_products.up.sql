-- Create products table
-- id, name, category_id, description, base_price, min_stock_level, barcode, is_active, created_at, updated_at
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    category_id UUID NOT NULL REFERENCES categories(id),
    description TEXT,
    base_price DECIMAL(10,2) NOT NULL,
    min_stock_level INT NOT NULL DEFAULT 0 CHECK (min_stock_level >= 0),
    barcode VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
-- idx_products_category_id on category_id
CREATE INDEX idx_products_category_id ON products(category_id);

-- Insert dummy data (5 products)
INSERT INTO products (name, category_id, description, base_price, min_stock_level, barcode, is_active) VALUES
('Paracetamol 500mg', (SELECT id FROM categories WHERE name = 'Medications' LIMIT 1), 'Pain reliever and fever reducer', 5000.00, 100, '1234567890123', true),
('Vitamin C 1000mg', (SELECT id FROM categories WHERE name = 'Vitamins & Supplements' LIMIT 1), 'Immune system support', 15000.00, 50, '1234567890124', true),
('Bandage Roll 5cm', (SELECT id FROM categories WHERE name = 'Medical Supplies' LIMIT 1), 'Elastic bandage for wound care', 8000.00, 200, '1234567890125', true),
('Face Mask 3-Ply', (SELECT id FROM categories WHERE name = 'Medical Supplies' LIMIT 1), 'Disposable face mask', 3000.00, 500, '1234567890126', true),
('Baby Shampoo 200ml', (SELECT id FROM categories WHERE name = 'Baby Care' LIMIT 1), 'Gentle baby shampoo', 25000.00, 30, '1234567890127', true);