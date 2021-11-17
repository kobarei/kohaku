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
    SELECT id
    FROM sora_node
    WHERE (
      (label = @label) AND
      (version = @version) AND
      (node_name = @node_name)
    )
);

-- name: InsertSoraConnection :exec
INSERT INTO sora_connection (
  timestamp,
  label, version, node_name,
  multistream, simulcast, spotlight,
  role, channel_id, session_id, client_id, connection_id
)
SELECT
  @timestamp,
  @label, @version, @node_name,
  @multistream, @simulcast, @spotlight,
  @role, @channel_id, @session_id, @client_id, @connection_id
WHERE
  NOT EXISTS (
    SELECT id
    FROM sora_connection
    WHERE (
      (channel_id = @channel_id) AND
      (session_id = @session_id) AND
      (client_id = @client_id) AND
      (connection_id = @connection_id)
    )
);

-- name: InsertErlangVmMemoryStats :exec
INSERT INTO erlang_vm_memory_stats (
  time,
  sora_label, sora_version, sora_node_name,
  stats_type,
  type_total, type_processes, type_processes_used, type_system,
  type_atom, type_atom_used, type_binary, type_code, type_ets
) VALUES (
  @time,
  @sora_version, @sora_label, @sora_node_name,
  @stats_type,
  @type_total, @type_processes, @type_processes_used, @type_system,
  @type_atom, @type_atom_used, @type_binary, @type_code, @type_ets
);