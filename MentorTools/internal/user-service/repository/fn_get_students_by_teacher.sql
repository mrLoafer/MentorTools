CREATE OR REPLACE FUNCTION fn_get_students_by_teacher(p_teacher_id INT)
    RETURNS TABLE (
                      id INT,
                      name VARCHAR,
                      email VARCHAR
                  ) LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
        SELECT s.id, s.name, s.email
        FROM user_links ul
                 JOIN users s ON ul.student_id = s.id
        WHERE ul.teacher_id = p_teacher_id;
END;
$$;

-- Permissions
ALTER FUNCTION fn_get_students_by_teacher(INT) OWNER TO user_service_user;
GRANT EXECUTE ON FUNCTION fn_get_students_by_teacher(INT) TO user_service_user;