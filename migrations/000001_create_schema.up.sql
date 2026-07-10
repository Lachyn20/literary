-- users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'admin',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- refresh_tokens
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- categories
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    parent_id UUID NULL REFERENCES categories(id) ON DELETE SET NULL
);

-- works
CREATE TYPE audience_type AS ENUM ('adult', 'children');

CREATE TABLE IF NOT EXISTS works (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    content TEXT,
    description TEXT,
    audience_type audience_type NOT NULL,
    publish_year INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_works_category_id ON works(category_id);
CREATE INDEX IF NOT EXISTS idx_works_publish_year ON works(publish_year);

-- full-text search vector for works (Turkmen "simple" configuration)
ALTER TABLE works ADD COLUMN IF NOT EXISTS search_vector tsvector;
CREATE INDEX IF NOT EXISTS idx_works_search_vector ON works USING GIN(search_vector);

CREATE FUNCTION works_search_vector_trigger() RETURNS trigger AS $$
begin
    new.search_vector := to_tsvector('simple', coalesce(new.title,'') || ' ' || coalesce(new.description,'') || ' ' || coalesce(new.content,''));
    return new;
end
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_works_search_vector ON works;
CREATE TRIGGER trg_works_search_vector
BEFORE INSERT OR UPDATE ON works
FOR EACH ROW EXECUTE FUNCTION works_search_vector_trigger();

-- translated_by_author
CREATE TABLE IF NOT EXISTS translated_by_author (
    id UUID PRIMARY KEY,
    original_author_name TEXT NOT NULL,
    original_language TEXT NOT NULL,
    work_title TEXT NOT NULL,
    notes TEXT
);

-- translated_into_languages
CREATE TABLE IF NOT EXISTS translated_into_languages (
    id UUID PRIMARY KEY,
    language_name TEXT NOT NULL,
    translator_name TEXT NOT NULL,
    work_title TEXT NOT NULL,
    notes TEXT
);

-- criticism_articles
CREATE TABLE IF NOT EXISTS criticism_articles (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author TEXT NOT NULL,
    publish_date DATE NOT NULL
);

-- books
CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    bibliographic_info TEXT,
    cover_image_path TEXT,
    pdf_path TEXT,
    page_count INT,
    published_year INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- book_photos
CREATE TABLE IF NOT EXISTS book_photos (
    id UUID PRIMARY KEY,
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    image_path TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_book_photos_book_id ON book_photos(book_id);

-- broadcasts
CREATE TYPE broadcast_type AS ENUM ('tv', 'radio');
CREATE TYPE file_type AS ENUM ('video', 'audio');

CREATE TABLE IF NOT EXISTS broadcasts (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    broadcast_type broadcast_type NOT NULL,
    channel_name TEXT NOT NULL,
    broadcast_date DATE NOT NULL,
    file_path TEXT,
    file_type file_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- theatre_productions
CREATE TABLE IF NOT EXISTS theatre_productions (
    id UUID PRIMARY KEY,
    play_title TEXT NOT NULL,
    theatre_name TEXT NOT NULL,
    premiere_date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- films
CREATE TYPE film_type AS ENUM ('film', 'animation');

CREATE TABLE IF NOT EXISTS films (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    film_type film_type NOT NULL,
    based_on_scenario BOOLEAN NOT NULL DEFAULT false,
    director TEXT,
    release_year INT,
    video_path TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- photo_archive
CREATE TYPE photo_category AS ENUM ('archive', 'personal');

CREATE TABLE IF NOT EXISTS photo_archive (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    image_path TEXT NOT NULL,
    description TEXT,
    taken_date DATE,
    category photo_category NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- biography
CREATE TABLE IF NOT EXISTS biography (
    id UUID PRIMARY KEY,
    content TEXT NOT NULL,
    photo_path TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- personal_letters
CREATE TABLE IF NOT EXISTS personal_letters (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    letter_date DATE NOT NULL,
    scan_image_path TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- external_links
CREATE TABLE IF NOT EXISTS external_links (
    id UUID PRIMARY KEY,
    site_name TEXT NOT NULL,
    url TEXT NOT NULL,
    category TEXT,
    notes TEXT
);
