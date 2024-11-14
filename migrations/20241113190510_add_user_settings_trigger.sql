-- +goose Up

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION create_user_settings()
    RETURNS TRIGGER AS $_$
BEGIN
    INSERT INTO user_settings (user_id)
    VALUES (NEW.id);
    RETURN NEW;
END;
$_$ LANGUAGE plpgsql;

CREATE TRIGGER after_user_insert
    AFTER INSERT ON users
    FOR EACH ROW
EXECUTE FUNCTION create_user_settings();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS after_user_insert ON users;
DROP FUNCTION IF EXISTS create_user_settings;
