CREATE TYPE audit_action AS ENUM ('created', 'updated', 'deleted', 'toggled');

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_id UUID REFERENCES flags(id) ON DELETE SET NULL,
    action audit_action NOT NULL,
    diff JSONB DEFAULT '{}'::jsonb,
    actor TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_audit_logs_flag_id ON audit_logs (flag_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs (created_at DESC);
