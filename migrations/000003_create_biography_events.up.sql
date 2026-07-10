-- Remove old content columns from biography (if they exist)
ALTER TABLE biography DROP COLUMN IF EXISTS content;
ALTER TABLE biography DROP COLUMN IF EXISTS content_tk;
ALTER TABLE biography DROP COLUMN IF EXISTS content_ru;
ALTER TABLE biography DROP COLUMN IF EXISTS content_en;

-- Create biography_events table
CREATE TABLE biography_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    biography_id UUID NOT NULL REFERENCES biography(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    title_tk VARCHAR(255),
    title_ru VARCHAR(255),
    title_en VARCHAR(255),
    description_tk TEXT,
    description_ru TEXT,
    description_en TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_biography_events_biography_id 
    ON biography_events(biography_id);
CREATE INDEX idx_biography_events_year 
    ON biography_events(year)