package kohaku

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Collector(c *gin.Context) {
	// TODO(v): validator 処理
	exporter := new(SoraStatsExporter)
	if err := c.Bind(exporter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// exporter.Type の conection.remote にのみ対応する
	// TODO(v): 将来的には conneciton.sora やそれ以外にも対応していく
	if exporter.Type == "connection.remote" {
		// context からコネクションプールを取り出す
		pool := c.MustGet("pool").(*pgxpool.Pool)
		if err := CollectorRemoteStats(pool, *exporter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
		return
	}

	c.Status(http.StatusBadRequest)
}
