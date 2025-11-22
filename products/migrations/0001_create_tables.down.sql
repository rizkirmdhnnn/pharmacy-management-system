-- Drop triggers
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (order matters due to foreign key constraints)
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;

