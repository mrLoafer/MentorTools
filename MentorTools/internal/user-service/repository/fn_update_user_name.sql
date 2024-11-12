CREATE OR REPLACE FUNCTION fn_update_user_name(
    p_user_id INT,
    p_name VARCHAR
) RETURNS TABLE (
                    code VARCHAR,
                    message VARCHAR
                ) LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET name = p_name WHERE id = p_user_id;

    IF NOT FOUND THEN
        RETURN QUERY SELECT 'USER001', 'User not found';
    ELSE
        RETURN QUERY SELECT 'SUCCESS', 'User name updated successfully';
    END IF;
END;
$$;

ALTER FUNCTION fn_update_user_name(INT, VARCHAR) OWNER TO auth_user;
GRANT EXECUTE ON FUNCTION fn_update_user_name(INT, VARCHAR) TO auth_user;
