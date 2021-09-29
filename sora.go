package kohaku

import (
	"encoding/json"
	"time"
)

// TODO: validator 処理の追加

// type は PeerConnection / SoraConnection
// type: connection.remote / type: connection.sora
type SoraStatsExporter struct {
	Type string `json:"type" binding:"required"`

	Label     string    `json:"label"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp" binding:"required"`

	// TODO: binding:"required,len=8" を追加する
	Role string `json:"role"`

	ChannelID    string `json:"channel_id" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	ConnectionID string `json:"connection_id" binding:"required,len=26"`

	// TODO(v): required にする
	Multistream bool `json:"multistream"`
	Simulcast   bool `json:"simulcast"`
	Spotlight   bool `json:"spotlight"`

	Stats []json.RawMessage `json:"stats" binding:"required"`
}
