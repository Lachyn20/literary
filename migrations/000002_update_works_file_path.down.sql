ALTER TABLE works DROP COLUMN file_path;
ALTER TABLE works ADD COLUMN content TEXT;

-- restore trigger function to include content again
CREATE OR REPLACE FUNCTION works_search_vector_trigger() RETURNS trigger AS $$
begin
	new.search_vector := to_tsvector('simple', coalesce(new.title,'') || ' ' || coalesce(new.description,'') || ' ' || coalesce(new.content,''));
	return new;
end
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_works_search_vector ON works;
CREATE TRIGGER trg_works_search_vector  
BEFORE INSERT OR UPDATE ON works
FOR EACH ROW EXECUTE FUNCTION works_search_vector_trigger();
