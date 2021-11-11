package kohaku

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type Server struct {
	config *KohakuConfig
	pool   *pgxpool.Pool
	http.Server
}

func NewServer(c *KohakuConfig, pool *pgxpool.Pool) *Server {
	r := gin.New()

	r.Use(httpLogger())
	r.Use(gin.Recovery())

	r.Use(validateHttpVersion())

	// TODO(v): こいつ自身の統計情報を /stats でとれた方がいい

	h2s := &http2.Server{
		MaxConcurrentStreams: c.Http2MaxConcurrentStreams,
		MaxReadFrameSize:     c.Http2MaxReadFrameSize,
		IdleTimeout:          time.Duration(c.Http2IdleTimeout) * time.Second,
	}

	s := &Server{
		config: c,
		pool:   pool,
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", c.CollectorPort),
			Handler: h2c.NewHandler(r, h2s),
		},
	}

	http2H2c := c.Http2H2c
	if !http2H2c {
		if c.Http2VerifyCacertPath != "" {
			clientCAPath := c.Http2VerifyCacertPath
			clientCA, err := ioutil.ReadFile(clientCAPath)
			if err != nil {
				panic(err)
			}

			certPool := x509.NewCertPool()
			ok := certPool.AppendCertsFromPEM(clientCA)
			if !ok {
				panic("failed to append certificates")
			}

			tlsConfig := &tls.Config{
				ClientAuth: tls.VerifyClientCertIfGiven,
				ClientCAs:  certPool,
			}
			s.Server.TLSConfig = tlsConfig
		}
	}

	// 統計情報を突っ込むところ
	r.POST("/collector", s.Collector)
	// ヘルスチェック
	r.POST("/status", s.Status)

	// Custom Validation Functions の登録
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// TODO: タグ名を変更する
		v.RegisterValidation("maxb", MaximumNumberOfBytesFunc)
	}

	return s
}

func (s *Server) Start(c *KohakuConfig) error {
	http2H2c := c.Http2H2c

	if http2H2c {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
	} else {
		http2FullchainFile := c.Http2FullchainFile
		http2PrivkeyFile := c.Http2PrivkeyFile
		if err := s.ListenAndServeTLS(http2FullchainFile, http2PrivkeyFile); err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

func validateHttpVersion() gin.HandlerFunc {
	return func(c *gin.Context) {

		// prior knowledge ではない場合
		if upgrade, ok := c.Request.Header["Upgrade"]; ok && upgrade[0] == "h2c" {
			return
		}

		if c.Request.Proto != "HTTP/2.0" {
			err := fmt.Errorf("UNSUPPORTED-HTTP-VERSION: %s", c.Request.Proto)
			// TODO: 505 を返すかの検討
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}

func httpLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		event := logEvent(c.Writer.Status())

		req := c.Request

		event.
			Int("status", c.Writer.Status()).
			Str("address", req.RemoteAddr).
			Str("method", req.Method).
			Str("path", req.URL.Path).
			Int64("len", req.ContentLength).
			Msg("")
	}
}

func logEvent(status int) *zerolog.Event {
	var event *zerolog.Event

	switch status / 100 {
	case 5:
		event = zlog.Error()
	case 4:
		event = zlog.Warn()
	default:
		event = zlog.Info()
	}

	return event
}
