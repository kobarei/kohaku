package kohaku

import (
	"encoding/json"
	"time"
)

// type は PeerConnection / SoraConnection
// type: peer-connection / type: sora-connection
type SoraStatsExporter struct {
	Type string `json:"type" validate:"required"`

	Label     string    `json:"label"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp" validate:"required"`

	ChannelID    string `json:"channel_id" validate:"required"`
	ClientID     string `json:"client_id" validate:"required"`
	ConnectionID string `json:"connection_id" validate:"required"`

	Stats []json.RawMessage `json:"stats" validate:"required"`
}
