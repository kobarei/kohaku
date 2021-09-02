package kohaku

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"
)

var (
	ConfigFilePath = flag.String("c", "./config.yaml", "kohaku 設定ファイルへのパス(yaml)")
	Config         *KohakuConfig
)

type KohakuConfig struct {
	Debug bool `yaml:"debug"`

	LogDir    string `yaml:"log_dir"`
	LogName   string `yaml:"log_name"`
	LogStdout bool   `yaml:"log_stdout"`

	CollectorPort int `yaml:"collector_port"`

	PostgresHost     string `yaml:"postgres_host"`
	PostgresPort     int    `yaml:"postgres_port"`
	PostgresUser     string `yaml:"postgres_user"`
	PostgresPassword string `yaml:"postgres_password"`
	PostgresDbname   string `yaml:"postgres_dbaname"`
	PostgresSslmode  string `yaml:"postgres_sslmode"`

	Http2CertFilePath         string `yaml:"http2_cert_file_path"`
	Http2KeyFilePath          string `yaml:"http2_key_file_path"`
	Http2H2c                  bool   `yaml:"http2_h2c"`
	Http2MaxConcurrentStreams int    `yaml:"http2_max_concurrent_streams"`
	Http2IdelTimeout          int    `yaml:"http2_idel_timeout"`
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

// LoadConfig 設定ファイルのパスからファイルを読み込み、設定値をバインドした AyameLaboConfig を返します
func LoadConfig(configPath string) (*KohakuConfig, error) {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config KohakuConfig
	if err := yaml.Unmarshal(buf, &config); err != nil {
		return nil, fmt.Errorf("AyameLaboConfig bind error: %s", err)
	}
	return &config, nil
}
