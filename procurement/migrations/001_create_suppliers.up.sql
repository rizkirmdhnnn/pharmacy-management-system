-- Create suppliers table
-- id, name, contact_person, email, phone, address, city, country, is_active, created_at, updated_at
CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    contact_person VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_suppliers_name ON suppliers(name);

-- Insert dummy data (5 suppliers)
INSERT INTO suppliers (name, contact_person, email, phone, address, city, country, is_active) VALUES
('PT. Medika Sejahtera', 'Budi Santoso', 'budi@medika.com', '081234567890', 'Jl. Sudirman No. 123', 'Jakarta', 'Indonesia', true),
('CV. Farmasi Nusantara', 'Siti Nurhaliza', 'siti@farmasi.com', '081987654321', 'Jl. Gatot Subroto No. 456', 'Bandung', 'Indonesia', true),
('PT. Kesehatan Prima', 'Ahmad Fauzi', 'ahmad@kesehatan.com', '082112345678', 'Jl. Thamrin No. 789', 'Surabaya', 'Indonesia', true),
('UD. Apotek Sehat', 'Dewi Sartika', 'dewi@apotek.com', '082234567890', 'Jl. Merdeka No. 321', 'Yogyakarta', 'Indonesia', true),
('PT. Obat Terpercaya', 'Rudi Hartono', 'rudi@obat.com', '083345678901', 'Jl. Diponegoro No. 654', 'Medan', 'Indonesia', true);

