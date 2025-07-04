-- Add indexes for better search performance
CREATE INDEX IF NOT EXISTS idx_books_title ON books (LOWER(title) text_pattern_ops);
CREATE INDEX IF NOT EXISTS idx_books_author ON books (LOWER(author) text_pattern_ops);
CREATE INDEX IF NOT EXISTS idx_books_year ON books (year);

-- Add a combined index for title and author if you often search both fields together
CREATE INDEX IF NOT EXISTS idx_books_title_author ON books (LOWER(title) text_pattern_ops, LOWER(author) text_pattern_ops);

-- Add a GIN index for full-text search
ALTER TABLE books ADD COLUMN IF NOT EXISTS search_vector tsvector;
UPDATE books SET search_vector = to_tsvector('english', title || ' ' || author || ' ' || COALESCE(year::text, ''));
CREATE INDEX IF NOT EXISTS idx_books_search ON books USING GIN(search_vector);

-- Create a trigger to update the search_vector column
CREATE OR REPLACE FUNCTION books_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector = to_tsvector('english', 
        COALESCE(NEW.title, '') || ' ' || 
        COALESCE(NEW.author, '') || ' ' || 
        COALESCE(NEW.year::text, '')
    );
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER books_search_vector_update_trigger
BEFORE INSERT OR UPDATE ON books
FOR EACH ROW EXECUTE FUNCTION books_search_vector_update();
