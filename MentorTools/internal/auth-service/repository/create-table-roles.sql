CREATE TABLE IF NOT EXISTS roles (
                                     id SERIAL PRIMARY KEY,
                                     role_name VARCHAR(50) UNIQUE NOT NULL,
                                     description VARCHAR(255)
);

COMMENT ON COLUMN public.roles.id IS 'Identifier';
COMMENT ON COLUMN public.roles.role_name IS 'Name of role';
COMMENT ON COLUMN public.roles.description IS 'Description';