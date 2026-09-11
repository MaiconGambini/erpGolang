CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name TEXT NOT NULL,
    document TEXT,
    email TEXT,
    phone TEXT,
    active BOOLEAN NOT NULL DEFAULT true,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_suppliers_tenant_id ON suppliers(tenant_id);
CREATE UNIQUE INDEX idx_suppliers_tenant_document_unique
    ON suppliers(tenant_id, document)
    WHERE document IS NOT NULL AND deleted_at IS NULL;
