-- Create purchases table
-- id, purchase_number, supplier_id, purchase_date, total_amount, status, notes, created_at, updated_at, created_by
CREATE TABLE purchases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    purchase_number VARCHAR(50) NOT NULL UNIQUE,
    supplier_id UUID NOT NULL REFERENCES suppliers(id),
    purchase_date DATE NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'cancelled')),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL
);

-- Create indexes
CREATE INDEX idx_purchases_supplier_id ON purchases(supplier_id);
CREATE INDEX idx_purchases_purchase_date ON purchases(purchase_date);
CREATE INDEX idx_purchases_status ON purchases(status);

-- Insert dummy data (5 purchases)
-- Note: created_by uses a dummy UUID, replace with actual user ID in production
INSERT INTO purchases (purchase_number, supplier_id, purchase_date, total_amount, status, notes, created_by) VALUES
('PO-001', (SELECT id FROM suppliers WHERE name = 'PT. Medika Sejahtera' LIMIT 1), CURRENT_DATE, 1500000.00, 'pending', 'Expected delivery: ' || (CURRENT_DATE + INTERVAL '7 days')::text, gen_random_uuid()),
('PO-002', (SELECT id FROM suppliers WHERE name = 'CV. Farmasi Nusantara' LIMIT 1), CURRENT_DATE - INTERVAL '5 days', 2400000.00, 'completed', 'Expected delivery: ' || (CURRENT_DATE + INTERVAL '2 days')::text, gen_random_uuid()),
('PO-003', (SELECT id FROM suppliers WHERE name = 'PT. Kesehatan Prima' LIMIT 1), CURRENT_DATE - INTERVAL '10 days', 800000.00, 'pending', 'Expected delivery: ' || (CURRENT_DATE + INTERVAL '10 days')::text, gen_random_uuid()),
('PO-004', (SELECT id FROM suppliers WHERE name = 'UD. Apotek Sehat' LIMIT 1), CURRENT_DATE - INTERVAL '3 days', 3200000.00, 'pending', 'Expected delivery: ' || (CURRENT_DATE + INTERVAL '5 days')::text, gen_random_uuid()),
('PO-005', (SELECT id FROM suppliers WHERE name = 'PT. Obat Terpercaya' LIMIT 1), CURRENT_DATE - INTERVAL '1 day', 1800000.00, 'completed', 'Expected delivery: ' || (CURRENT_DATE + INTERVAL '3 days')::text, gen_random_uuid());

