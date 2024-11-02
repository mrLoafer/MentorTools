DROP FUNCTION IF EXISTS public.fn_create_user(varchar, varchar, varchar, varchar);

CREATE OR REPLACE FUNCTION public.fn_create_user(
    p_username character varying,
    p_password_hash character varying,
    p_email character varying,
    p_role_name character varying,
    OUT code character varying,
    OUT message character varying
)
    LANGUAGE plpgsql
AS $function$
DECLARE
    v_role_id int;
BEGIN
    -- Retrieve role_id based on role_name
    SELECT id INTO v_role_id FROM roles WHERE role_name = p_role_name;

    -- If role_id is null, return a custom error with code AUTH0003
    IF v_role_id IS NULL THEN
        code := 'AUTH0003';
        message := 'Role does not exist';
        RETURN;
    END IF;

    -- Attempt to insert the new user; if email already exists, return error
    INSERT INTO users (username, password_hash, email, role_id)
    VALUES (p_username, p_password_hash, p_email, v_role_id)
    ON CONFLICT (email) DO NOTHING;

    -- If no rows were inserted, return custom error with code AUTH0001
    IF NOT FOUND THEN
        code := 'AUTH0001';
        message := 'User with the same email already exists';
        RETURN;
    END IF;

    -- If successful, set success code and message
    code := 'SUCCESS';
    message := 'User created successfully';
END;
$function$;

-- Permissions
    ALTER FUNCTION public.fn_create_user(varchar, varchar, varchar, varchar) OWNER TO auth_user;
GRANT ALL ON FUNCTION public.fn_create_user(varchar, varchar, varchar, varchar) TO auth_user;