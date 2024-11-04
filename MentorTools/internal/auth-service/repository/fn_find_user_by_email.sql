CREATE OR REPLACE FUNCTION public.fn_find_user_by_email(p_email VARCHAR)
    RETURNS TABLE (
                      code VARCHAR,
                      message VARCHAR,
                      user_id INT,
                      email VARCHAR,
                      password_hash VARCHAR,
                      role_name VARCHAR
                  )
    LANGUAGE plpgsql
AS $$
BEGIN
    -- Check if user with the specified email exists
    RETURN QUERY
        SELECT
            'SUCCESS'::VARCHAR AS code,
            'User found'::VARCHAR AS message,
            u.id::INT AS user_id,
            u.email::VARCHAR AS email,
            u.password_hash::VARCHAR AS password_hash,
            r.role_name::VARCHAR AS role_name
        FROM users u
                 LEFT JOIN roles r ON r.id = u.role_id
        WHERE u.email = p_email;

    -- If user is not found, return custom error response
    IF NOT FOUND THEN
        RETURN QUERY SELECT
                         'AUTH0005'::VARCHAR AS code,
                         'User not found'::VARCHAR AS message,
                         NULL::INT AS user_id,
                         NULL::VARCHAR AS email,
                         NULL::VARCHAR AS password_hash,
                         NULL::VARCHAR AS role_name;
    END IF;
END;
$$;

-- Permissions
    ALTER FUNCTION public.fn_find_user_by_email(varchar) OWNER TO auth_user;
GRANT ALL ON FUNCTION public.fn_find_user_by_email(varchar) TO auth_user;