CREATE TABLE public.users (
                              id serial4 NOT NULL, -- Identifier
                              username varchar(255) NOT NULL, -- Name of user (Name, Surname ..)
                              password_hash varchar(255) NOT NULL, -- Hash of password
                              email varchar(100) NOT NULL, -- email as user login
                              role_id INT NOT NULL, -- Reference to user role
                              created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
                              updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
                              CONSTRAINT users_email_key UNIQUE (email),
                              CONSTRAINT users_pkey PRIMARY KEY (id),
                              CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id) REFERENCES roles(id) -- Foreign key constraint for role
);

-- Column comments

COMMENT ON COLUMN public.users.id IS 'Identifier';
COMMENT ON COLUMN public.users.username IS 'Name of user (Name, Surname ..)';
COMMENT ON COLUMN public.users.password_hash IS 'Hash of password';
COMMENT ON COLUMN public.users.email IS 'Email as user login';
COMMENT ON COLUMN public.users.role_id IS 'Reference to user role';
COMMENT ON COLUMN public.users.created_at IS 'Timestamp when the user was created';
COMMENT ON COLUMN public.users.updated_at IS 'Timestamp when the user was last updated';