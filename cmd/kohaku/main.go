package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shiguredo/kohaku"

	zlog "github.com/rs/zerolog/log"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// curl -v --http2-prior-knowledge http://localhost:8080

// TODO(v): 特定の IP アドレスからしか受け付けないようにする
// TODO(v): 性能測定
// TODO(v): YAML 化

func init() {
	// 設定の読み込みと、Logger の準備
	flag.Parse()
	// load config
	if err := kohaku.LoadConfigFromFlags(kohaku.ConfigFilePath); err != nil {
		log.Fatalf("config file read error: %s", err)
	}
	err := kohaku.InitLogger(kohaku.Config.LogDir, kohaku.Config.LogName, kohaku.Config.Debug, kohaku.Config.LogStdout)
	if err != nil {
		log.Fatalf("logger building failed. %s", err)
	}
	zlog.Info().
		Msg("FinishInitProcess")
}

func pgxPoolMiddleware(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("pool", pool)
		c.Next()
	}
}

func NewDB(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		// TODO(v): エラーメッセージを修正する
		fmt.Fprintf(os.Stderr, "Unable to parse url: %v\n", err)
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		// TODO(v): エラーメッセージを修正する
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		// TODO(v): エラーメッセージを修正する
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	return pool, nil
}

type Server struct {
	config *kohaku.KohakuConfig
	pool   *pgxpool.Pool
	http.Server
}

func NewServer(c *kohaku.KohakuConfig, pool *pgxpool.Pool) *Server {
	r := gin.New()

	// TODO(v): カスタムコンテキストに Pool を渡すかたちでいいのかどうか確認する
	r.Use(pgxPoolMiddleware(pool))

	// TODO(v): zerolog に切り替える
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 統計情報を突っ込むところ
	r.POST("/collector", kohaku.Collector)

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

	return s
}

func (s *Server) Start(c *kohaku.KohakuConfig) error {
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

func main() {
	var connStr = kohaku.Config.PostgresURL
	pool, err := NewDB(context.Background(), connStr)
	if err != nil {
		// TODO: 共通化できるのであればエラーメッセージはここで出力する
		os.Exit(1)
	}
	defer pool.Close()

	s := NewServer(kohaku.Config, pool)

	if err := s.Start(kohaku.Config); err != nil {
		log.Fatal(err)
	}
}
