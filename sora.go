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

	Role string `json:"role" binding:"required,len=8"`

	ChannelID    string `json:"channel_id" binding:"required"`
	SessionID    string `json:"session_id" binding:"required,len=26"`
	ClientID     string `json:"client_id" binding:"required"`
	ConnectionID string `json:"connection_id" binding:"required,len=26"`

	Multistream bool `json:"multistream" binding:"required"`
	Simulcast   bool `json:"simulcast" binding:"required"`
	Spotlight   bool `json:"spotlight" binding:"required"`

	Stats []json.RawMessage `json:"stats" binding:"required"`
}
