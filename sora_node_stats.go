package kohaku

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	db "github.com/shiguredo/kohaku/db/sqlc"
)

func toNumeric(n uint64) pgtype.Numeric {
	var num pgtype.Numeric
	num.Set(n)
	return num
}

// TODO(v): sqlc 化
func (s *Server) CollectorSoraNodeErlangVmStats(c *gin.Context, stats SoraNodeErlangVmStats) error {
	if err := s.InsertSoraNode(c, stats); err != nil {
		return err
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

			if err := s.query.InsertErlangVmMemoryStats(c, db.InsertErlangVmMemoryStatsParams{
				Time:              stats.Timestamp,
				SoraVersion:       stats.Version,
				SoraLabel:         stats.Label,
				SoraNodeName:      stats.NodeName,
				StatsType:         e.Type,
				TypeTotal:         toNumeric(e.Total),
				TypeProcesses:     toNumeric(e.Processes),
				TypeProcessesUsed: toNumeric(e.ProcessesUsed),
				TypeSystem:        toNumeric(e.System),
				TypeAtom:          toNumeric(e.Atom),
				TypeAtomUsed:      toNumeric(e.AtomUsed),
				TypeBinary:        toNumeric(e.Binary),
				TypeCode:          toNumeric(e.Code),
				TypeEts:           toNumeric(e.ETS),
			}); err != nil {
				return err
			}
		default:
			// TODO: return err にする
			fmt.Println(erlangVmStats.Type)
		}
	}
	return nil
}

func (s *Server) InsertSoraNode(c *gin.Context, stats SoraNodeErlangVmStats) error {
	if err := s.query.InsertSoraNode(c, db.InsertSoraNodeParams{
		Timestamp: stats.Timestamp,
		Label:     stats.Label,
		Version:   stats.Version,
		NodeName:  stats.NodeName,
	}); err != nil {
		return err
	}
	return nil
}
