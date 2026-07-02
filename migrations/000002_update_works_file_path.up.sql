-- Update trigger function to stop referencing `content` before dropping the column
CREATE OR REPLACE FUNCTION works_search_vector_trigger() RETURNS trigger AS $$
begin
	new.search_vector := to_tsvector('simple', coalesce(new.title,'') || ' ' || coalesce(new.description,''));
	return new;
end
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_works_search_vector ON works;
CREATE TRIGGER trg_works_search_vector
BEFORE INSERT OR UPDATE ON works
FOR EACH ROW EXECUTE FUNCTION works_search_vector_trigger();

ALTER TABLE works ADD COLUMN file_path TEXT;
ALTER TABLE works DROP COLUMN content;
