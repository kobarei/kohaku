-- name: InsertSoraNode :exec
INSERT INTO sora_node (
  timestamp,
  label, version, node_name
)
SELECT
  @timestamp,
  @label, @version, @node_name
WHERE
  NOT EXISTS (
    SELECT label
    FROM sora_node
    WHERE (
      (label = @label) AND
      (node_name = @node_name) AND
      (version = @version)
    )
);