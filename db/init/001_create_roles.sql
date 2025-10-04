DO $$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'news_user') THEN
            CREATE ROLE news_user WITH LOGIN PASSWORD 'news_user' CREATEDB;
            GRANT ALL PRIVILEGES ON DATABASE news_db TO news_user;
        END IF;
    END
$$;