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