CREATE OR REPLACE FUNCTION fn_create_link(
    p_teacher_id INT,
    p_student_email VARCHAR
) RETURNS TABLE (
                    code VARCHAR,
                    message VARCHAR
                ) LANGUAGE plpgsql AS $$
DECLARE
    v_student_id INT;
BEGIN
    -- Insert link with student_id lookup and conflict handling in a single query
    INSERT INTO user_links (teacher_id, student_id)
    SELECT p_teacher_id, id
    FROM users
    WHERE email = p_student_email
    ON CONFLICT (teacher_id, student_id) DO NOTHING;

    -- Check if the insert succeeded; if not, return a general error message
    IF NOT FOUND THEN
        RETURN QUERY SELECT 'LINK005', 'Link creation failed or link already exists';
    END IF;

    -- Return success message
    RETURN QUERY SELECT 'SUCCESS', 'Link created successfully';
END;
$$;

-- Permissions
ALTER FUNCTION fn_create_link(INT, VARCHAR) OWNER TO auth_user;
GRANT EXECUTE ON FUNCTION fn_create_link(INT, VARCHAR) TO auth_user;