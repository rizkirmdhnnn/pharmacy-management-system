-- Create categories table
-- id, name, description, created_at, updated_at
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_categories_name ON categories(name);

-- Insert dummy data (5 categories)
INSERT INTO categories (name, description) VALUES
('Medications', 'Prescription and over-the-counter medications'),
('Vitamins & Supplements', 'Vitamins, minerals, and dietary supplements'),
('Medical Supplies', 'Bandages, syringes, and medical equipment'),
('Personal Care', 'Skincare, hygiene, and personal care products'),
('Baby Care', 'Baby products and infant care items');