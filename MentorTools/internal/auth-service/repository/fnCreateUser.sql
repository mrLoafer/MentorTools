-- DROP FUNCTION public."fnCreateUser"(varchar, varchar, varchar, int);

CREATE OR REPLACE FUNCTION public."fnCreateUser"(
    username character varying,
    password_hash character varying,
    email character varying,
    role_id int
)
    RETURNS void
    LANGUAGE plpgsql
AS $function$
BEGIN
    -- Attempt to insert the new user; raise error if email already exists
    INSERT INTO users (username, password_hash, email, role_id)
    SELECT username, password_hash, email, role_id
    WHERE NOT EXISTS (
        SELECT 1 FROM users WHERE email = email
    );

    -- If no rows were inserted, raise a custom error with code AUTH0001
    IF NOT FOUND THEN
        RAISE EXCEPTION 'AUTH0001' USING ERRCODE = 'P0001';
    END IF;
END;
$function$
;

-- Permissions

ALTER FUNCTION public."fnCreateUser"(varchar, varchar, varchar, int) OWNER TO auth_user;
GRANT ALL ON FUNCTION public."fnCreateUser"(varchar, varchar, varchar, int) TO auth_user;