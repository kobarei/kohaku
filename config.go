package kohaku

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
)

var (
	ConfigFilePath = flag.String("c", "./config.toml", "kohaku 設定ファイルへのパス(toml)")
	Config         *KohakuConfig
)

type KohakuConfig struct {
	LogDebug  bool   `toml:"log_debug"`
	LogDir    string `toml:"log_dir"`
	LogName   string `toml:"log_name"`
	LogStdout bool   `toml:"log_stdout"`

	CollectorPort int `toml:"collector_port"`

	TimescaleURL          string `toml:"timescale_url"`
	TimescaleSSLMode      string `toml:"timescale_sslmode"`
	TimescaleRootcertFile string `toml:"timescale_rootcert_file"`

	// TODO(v): 名前検討
	HTTP2FullchainFile string `toml:"http2_fullchain_file"`
	// TODO(v): 名前検討
	HTTP2PrivkeyFile string `toml:"http2_privkey_file"`
	// TODO: 名前検討
	HTTP2VerifyCacertPath string `toml:"http2_verify_cacert_path"`

	HTTP2H2c                  bool   `toml:"http2_h2c"`
	HTTP2MaxConcurrentStreams uint32 `toml:"http2_max_concurrent_streams"`
	HTTP2MaxReadFrameSize     uint32 `toml:"http2_max_read_frame_size"`
	HTTP2IdleTimeout          uint32 `toml:"http2_idle_timeout"`
}

// LoadConfigFromFlags 起動パラメータから設定ファイルを読み込みます
func LoadConfigFromFlags(configPath *string) error {
	tmpConfig, err := LoadConfig(*configPath)
	log.Printf("config file path: %s", *configPath)
	if err != nil {
		return err
	}
	Config = tmpConfig

	return nil
}

// LoadConfig 設定ファイルのパスからファイルを読み込み、設定値をバインドした KohakuConfig を返します
func LoadConfig(configPath string) (*KohakuConfig, error) {
	buf, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config KohakuConfig
	if err := toml.Unmarshal(buf, &config); err != nil {
		return nil, fmt.Errorf("KohakuConfig bind error: %s", err)
	}
	return &config, nil
}
