CREATE TABLE events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    banner_id uuid NOT NULL REFERENCES banners(id) ON DELETE CASCADE,
    type varchar(20) not null check (type in ('click', 'impression')),
    ip inet,
    user_agent text,
    created_at timestamptz DEFAULT now()
);
