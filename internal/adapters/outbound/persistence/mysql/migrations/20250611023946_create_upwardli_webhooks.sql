-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS upwardli.webhooks (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    webhook_name VARCHAR(255) NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    partner_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status VARCHAR(255) NOT NULL,
    failures INT DEFAULT 0,
    last_failure TIMESTAMP,
    deleted boolean DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS upwardli.webhooks;
-- +goose StatementEnd