-- name: CreateActivityLog :exec
INSERT INTO activity_log (user_id, action, resource_type, resource_id, resource_title, description)
VALUES (?, ?, ?, ?, ?, ?);

-- name: ListActivityLogs :many
SELECT * FROM activity_log
WHERE
    (CAST(@filter_action AS TEXT) = '' OR action = @filter_action)
    AND (CAST(@filter_search AS TEXT) = '' OR description LIKE '%' || @filter_search || '%')
ORDER BY created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountActivityLogs :one
SELECT COUNT(*) FROM activity_log
WHERE
    (CAST(@filter_action AS TEXT) = '' OR action = @filter_action)
    AND (CAST(@filter_search AS TEXT) = '' OR description LIKE '%' || @filter_search || '%');
