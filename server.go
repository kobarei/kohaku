package kohaku

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	config *KohakuConfig
	pool   *pgxpool.Pool
	http.Server
}

func NewServer(c *KohakuConfig, pool *pgxpool.Pool) *Server {
	r := gin.New()

	// TODO(v): zerolog に切り替える
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// TODO(v): ヘルスチェック用の /status みたいなのあった方がいい
	// TODO(v): こいつ自身の統計情報を /stats でとれた方がいい

	// TODO(v): HTTP/2 の設定は YAML へ
	h2s := &http2.Server{
		MaxConcurrentStreams: 250,
		MaxReadFrameSize:     1048576,
		IdleTimeout:          600 * time.Second,
	}

	s := &Server{
		config: c,
		pool:   pool,
		Server: http.Server{
			Addr: fmt.Sprintf(":%d", c.CollectorPort),
			// TODO(v): YAML で h2c と h2 を切り替えられるようにする
			Handler: h2c.NewHandler(r, h2s),
		},
	}

	// 統計情報を突っ込むところ
	r.POST("/collector", s.Collector)

	return s
}

func (s *Server) Start(c *KohakuConfig) error {
	http2H2c := c.Http2H2c

	if http2H2c {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
	} else {
		http2CertFilePath := c.Http2CertFilePath
		http2KeyFilePath := c.Http2KeyFilePath
		if err := s.ListenAndServeTLS(http2CertFilePath, http2KeyFilePath); err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}
