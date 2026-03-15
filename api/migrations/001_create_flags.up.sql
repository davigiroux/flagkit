CREATE TYPE flag_environment AS ENUM ('production', 'staging', 'development');

CREATE TABLE flags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    enabled BOOLEAN DEFAULT false,
    environment flag_environment NOT NULL DEFAULT 'development',
    rules JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_flags_key ON flags (key);
CREATE INDEX idx_flags_environment ON flags (environment);
