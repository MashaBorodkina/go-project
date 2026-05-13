CREATE TABLE banners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    image_url VARCHAR(500) NOT NULL CHECK (length(image_url) >= 1 AND length(image_url) <= 500),
    title varchar(200) check (length(title) >= 1 AND length(title) <= 200),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);