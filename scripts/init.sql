DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = 'messages') THEN
        CREATE TABLE messages (
            id SERIAL PRIMARY KEY,
            content TEXT NOT NULL,
            processed BOOLEAN DEFAULT FALSE
        );
    END IF;
END $$;
