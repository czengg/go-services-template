-- name: CreateUpwardliWebhook :exec
INSERT INTO upwardli.webhooks (
        id,
        webhook_name,
        endpoint,
        partner_id,
        status,
        failures,
        last_failure
    )
VALUES (?, ?, ?, ?, ?, ?, ?);
-- name: GetUpwardliWebhookById :one
SELECT id,
    webhook_name,
    endpoint,
    partner_id,
    status,
    failures,
    last_failure,
    created_at,
    updated_at,
    deleted
FROM upwardli.webhooks
WHERE id = ?
    AND deleted = FALSE;
-- name: GetAllUpwardliWebhooks :many
SELECT id,
    webhook_name,
    endpoint,
    partner_id,
    status,
    failures,
    last_failure,
    created_at,
    updated_at,
    deleted
FROM upwardli.webhooks
WHERE deleted = FALSE
ORDER BY created_at DESC;
-- name: SoftDeleteUpwardliWebhook :exec
UPDATE upwardli.webhooks
SET deleted = TRUE,
    deleted_at = NOW(),
    updated_at = NOW()
WHERE id = ?
    AND deleted = FALSE;