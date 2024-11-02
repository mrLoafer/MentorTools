CREATE OR REPLACE FUNCTION public.fn_find_user_by_email(p_email VARCHAR)
    RETURNS RECORD
    LANGUAGE plpgsql
AS $$
DECLARE
    result RECORD;
BEGIN
    -- Check if user with the specified email exists
    SELECT
        'SUCCESS' AS code,
        'User found' AS message,
        u.id AS user_id,
        u.email,
        u.password_hash,
        r.role_name
    INTO result
    FROM users u
             LEFT JOIN roles r ON r.id = u.role_id
    WHERE u.email = p_email;

    -- If user is not found, set custom error response
    IF NOT FOUND THEN
        result := (SELECT 'AUTH0005' AS code, 'User not found' AS message, NULL, NULL, NULL, NULL);
    END IF;

    RETURN result;
END;
$$;

-- Permissions
    ALTER FUNCTION public.fn_find_user_by_email(varchar) OWNER TO auth_user;
GRANT ALL ON FUNCTION public.fn_find_user_by_email(varchar) TO auth_user;