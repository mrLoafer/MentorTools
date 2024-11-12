CREATE OR REPLACE FUNCTION fn_create_user_link(
    p_teacher_id INT,
    p_student_id INT
) RETURNS TABLE (
                    code VARCHAR,
                    message VARCHAR
                ) LANGUAGE plpgsql AS $$
BEGIN
    -- Попытка вставки новой связи, выполняя все проверки одним запросом
    WITH teacher_check AS (
        SELECT id
        FROM users
        WHERE id = p_teacher_id AND role = 'teacher'
    ),
         student_check AS (
             SELECT id
             FROM users
             WHERE id = p_student_id
         )
    INSERT INTO user_links (teacher_id, student_id)
    SELECT t.id, s.id
    FROM teacher_check t, student_check s
    WHERE NOT EXISTS (
        SELECT 1
        FROM user_links
        WHERE teacher_id = p_teacher_id
          AND student_id = p_student_id
    )
    RETURNING 'SUCCESS' AS code, 'Link created successfully' AS message;

    -- Проверка на отсутствие учителя
    IF NOT FOUND THEN
        IF NOT EXISTS (SELECT 1 FROM users WHERE id = p_teacher_id AND role = 'teacher') THEN
            RETURN QUERY SELECT 'USER002' AS code, 'User is not a teacher' AS message;
        ELSIF NOT EXISTS (SELECT 1 FROM users WHERE id = p_student_id) THEN
            RETURN QUERY SELECT 'USER003' AS code, 'Student not found' AS message;
        ELSE
            RETURN QUERY SELECT 'USER004' AS code, 'Link already exists' AS message;
        END IF;
    END IF;
END;
$$;

-- Назначение прав
    ALTER FUNCTION fn_create_user_link(INT, INT) OWNER TO auth_user;
GRANT EXECUTE ON FUNCTION fn_create_user_link(INT, INT) TO auth_user;