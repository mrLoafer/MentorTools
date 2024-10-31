DO
$$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'auth_user') THEN
            CREATE ROLE auth_user WITH LOGIN PASSWORD 'auth_password';
        END IF;
    END
$$;

DO
$$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'auth_db') THEN
            CREATE DATABASE auth_db OWNER auth_user;
        END IF;
    END
$$;

GRANT ALL PRIVILEGES ON DATABASE auth_db TO auth_user;