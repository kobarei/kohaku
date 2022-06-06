package kohaku

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	zlog "github.com/rs/zerolog/log"
)

// TODO: ログレベル、ログメッセージを変更する
func (s *Server) collector(c echo.Context) error {
	t := c.Request().Header.Get("x-sora-stats-exporter-type")
	switch t {
	case "connection.user-agent":
		// TODO(v): validator 処理
		stats := new(soraConnectionStats)
		if err := c.Bind(stats); err != nil {
			zlog.Debug().Str("type", t).Err(err).Msg("")
			return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		if err := s.collectorUserAgentStats(c, *stats); err != nil {
			zlog.Warn().Str("type", t).Err(err).Msg("")
			return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return c.NoContent(http.StatusNoContent)
	case "node.erlang-vm":
		stats := new(soraNodeErlangVMStats)
		if err := c.Bind(stats); err != nil {
			zlog.Debug().Str("type", t).Err(err).Msg("")
			return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		if err := s.collectorSoraNodeErlangVMStats(c, *stats); err != nil {
			zlog.Warn().Str("type", t).Err(err).Msg("")
			return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return c.NoContent(http.StatusNoContent)
	default:
		zlog.Warn().Str("type", t).Msgf("UNEXPECTED-TYPE")
		return c.NoContent(http.StatusBadRequest)
	}
}
