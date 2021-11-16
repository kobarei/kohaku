package kohaku

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	db "github.com/shiguredo/kohaku/db/sqlc"
)

// TODO(v): sqlc 化
func (s *Server) CollectorSoraNodeErlangVmStats(c *gin.Context, stats SoraNodeErlangVmStats) error {
	if err := s.InsertSoraNode(c, stats); err != nil {
		return err
	}

	erlangVm := &ErlangVm{
		Time: stats.Timestamp,

		Label:    stats.Label,
		Version:  stats.Version,
		NodeName: stats.NodeName,
	}

	for _, v := range stats.Stats {
		erlangVmStats := new(ErlangVmStats)
		if err := json.Unmarshal(v, &erlangVmStats); err != nil {
			return err
		}

		// type をみて struct をさらに別途デコードする
		switch erlangVmStats.Type {
		case "erlang-vm-memory":
			e := new(ErlangVmMemoryStats)
			if err := json.Unmarshal(v, &e); err != nil {
				return err
			}

			ds := goqu.Insert("erlang_vm_memory_stats").Rows(
				ErlangVmMemory{
					ErlangVm:            *erlangVm,
					ErlangVmMemoryStats: *e,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := s.pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		default:
			// TODO: return err にする
			fmt.Println(erlangVmStats.Type)
		}
	}
	return nil
}

func (s *Server) InsertSoraNode(ctx context.Context, stats SoraNodeErlangVmStats) error {
	// TODO: db.New 毎回していいのか？
	if err := s.query.InsertSoraNode(ctx, db.InsertSoraNodeParams{
		Timestamp: *stats.Timestamp,
		Label:     stats.Label,
		Version:   stats.Version,
		NodeName:  stats.NodeName,
	}); err != nil {
		return err
	}
	return nil
}
