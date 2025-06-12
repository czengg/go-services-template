-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS upwardli.consumers (
    id VARCHAR(255) NOT NULL,
    pcid VARCHAR(255) NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL,
    kyc_status VARCHAR(255) NOT NULL,
    tax_id_type VARCHAR(255) NOT NULL,
    tax_identifier VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS upwardli.consumers;
-- +goose StatementEnd
