package kohaku

import (
	"encoding/json"
	"time"
)

// type は PeerConnection / SoraConnection
// type: connection.remote / type: connection.sora
type SoraStatsExporter struct {
	Type string `json:"type" validate:"required"`

	Label     string    `json:"label"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp" validate:"required"`

	ChannelID    string `json:"channel_id" validate:"required"`
	ClientID     string `json:"client_id" validate:"required"`
	ConnectionID string `json:"connection_id" validate:"required"`

	// TODO(v): multistream や simulcast などを突っ込む

	Stats []json.RawMessage `json:"stats" validate:"required"`
}
