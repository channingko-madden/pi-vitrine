-- +goose Up
-- +goose StatementBegin
DO 
$do$
BEGIN
IF NOT EXISTS (SELECT column_name FROM information_schema.columns WHERE table_name = 'devices' AND column_name = 'location') THEN
        ALTER TABLE devices ADD location VARCHAR(255) NOT NULL DEFAULT '';
END IF;
END
$do$
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO 
$do$
BEGIN
IF EXISTS (SELECT column_name FROM information_schema.columns WHERE table_name = 'devices' AND column_name = 'location') THEN
        ALTER TABLE devices DROP COLUMN location;
END IF;
END
$do$
-- +goose StatementEnd
