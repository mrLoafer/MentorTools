-- Проверка существования базы данных и её создание
\connect postgres
SELECT 'CREATE DATABASE auth_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'auth_db')\gexec

-- Проверка существования пользователя и его создание
DO $$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'auth_user') THEN
            CREATE USER auth_user WITH PASSWORD 'auth_password';
        END IF;
    END $$;

-- Назначение владельца для базы данных
ALTER DATABASE auth_db OWNER TO auth_user;

-- Подключение к созданной базе данных
\connect auth_db

-- Назначение владельца для схемы public и привилегий на уровне схемы
ALTER SCHEMA public OWNER TO auth_user;
GRANT ALL PRIVILEGES ON SCHEMA public TO auth_user;

-- Установка привилегий на создание таблиц по умолчанию
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO auth_user;