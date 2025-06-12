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
-- name: SaveUpwardliConsumer :exec
INSERT INTO upwardli.consumers (
        id,
        pcid,
        external_id,
        is_active,
        kyc_status,
        tax_id_type
    )
VALUES (
        ?,
        ?,
        ?,
        ?,
        ?,
        ?
    ) ON DUPLICATE KEY
UPDATE pcid =
VALUES(pcid),
    is_active =
VALUES(is_active),
    kyc_status =
VALUES(kyc_status),
    tax_id_type =
VALUES(tax_id_type)