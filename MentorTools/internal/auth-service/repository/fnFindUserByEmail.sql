DROP FUNCTION IF EXISTS public.fnFindUserByEmail(varchar);

CREATE OR REPLACE FUNCTION public.fnFindUserByEmail(p_email VARCHAR)
    RETURNS TABLE(email VARCHAR, password_hash VARCHAR, username VARCHAR, role_name VARCHAR) AS $$
BEGIN
    RETURN QUERY
        SELECT
            u.id,
            u.email,
            u.password_hash,
            u.username,
            r.role_name
        FROM users u
                 LEFT JOIN roles r ON r.id = u.role_id
        WHERE u.email = p_email;
END;
$$ LANGUAGE plpgsql;

-- Permissions

ALTER FUNCTION public.fnFindUserByEmail(varchar) OWNER TO auth_user;
GRANT ALL ON FUNCTION public.fnFindUserByEmail(varchar) TO auth_user;
