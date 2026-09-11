ALTER TABLE customers
    ADD COLUMN document_type TEXT CHECK (document_type IN ('cpf', 'cnpj')),
    ADD COLUMN postal_code TEXT,
    ADD COLUMN street TEXT,
    ADD COLUMN street_number TEXT,
    ADD COLUMN city TEXT,
    ADD COLUMN state TEXT CHECK (state IS NULL OR length(state) = 2);

ALTER TABLE suppliers
    ADD COLUMN document_type TEXT CHECK (document_type IN ('cpf', 'cnpj')),
    ADD COLUMN postal_code TEXT,
    ADD COLUMN street TEXT,
    ADD COLUMN street_number TEXT,
    ADD COLUMN city TEXT,
    ADD COLUMN state TEXT CHECK (state IS NULL OR length(state) = 2);

ALTER TABLE products
    ADD COLUMN unit TEXT NOT NULL DEFAULT 'UN',
    ADD COLUMN barcode TEXT;
