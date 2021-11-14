package kohaku

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TODO(v): sqlc 化
func CollectorSoraNodeErlangVmMemoryStats(pool *pgxpool.Pool, stats SoraNodeErlangVmMemoryStats) error {
	if err := InsertSoraNodeErlangVmMemoryStats(context.Background(), pool, stats); err != nil {
		return err
	}
	ds := goqu.Insert("sora_node_erlang_vm_stats").Rows(
		SoraNodeErlangVM{
			TotalMemory:        stats.TotalMemory,
			TotalProcesses:     stats.TotalProcesses,
			TotalProcessesUsed: stats.TotalProcessesUsed,
			TotalSystem:        stats.TotalSystem,
			TotalAtom:          stats.TotalAtom,
			TotalAtomUsed:      stats.TotalAtomUsed,
			TotalBinary:        stats.TotalBinary,
			TotalCode:          stats.TotalCode,
			TotalETS:           stats.TotalETS,
		},
	)
	insertSQL, _, _ := ds.ToSQL()
	_, err := pool.Exec(context.Background(), insertSQL)
	if err != nil {
		return err
	}

	return nil
}

// TODO(v): sqlc 化
func InsertSoraNodeErlangVmMemoryStats(ctx context.Context, pool *pgxpool.Pool, stats SoraNodeErlangVmMemoryStats) error {
	sq := goqu.Select("channel_id").
		From("sora_node").
		Where(goqu.Ex{
			"label":     stats.Label,
			"node_name": stats.NodeName,
			"version":   stats.Version,
		})
	le := goqu.L("NOT EXISTS ?", sq)

	ds := goqu.Insert("sora_node").
		Cols(
			"timestamp",

			"label",
			"version",
			"node_name",
		).
		FromQuery(
			goqu.Select(
				goqu.L("?, ?, ?, ?",
					stats.Timestamp,

					stats.Label,
					stats.Version,
					stats.NodeName,
				),
			).Where(le))
	insertSQL, _, _ := ds.ToSQL()
	if _, err := pool.Exec(ctx, insertSQL); err != nil {
		return err
	}

	return nil
}
