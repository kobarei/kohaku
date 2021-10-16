package kohaku

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

// TODO: validator 処理の追加

// type は PeerConnection / SoraConnection
// type: connection.remote / type: connection.sora
type SoraStatsExporter struct {
	Type string `json:"type" binding:"required"`

	Label     string     `json:"label"`
	Version   string     `json:"version"`
	Timestamp *time.Time `json:"timestamp" binding:"required"`

	// TODO(v): NodeName string `json:"node_name"`

	Role string `json:"role" binding:"required,len=8"`

	ChannelID    string `json:"channel_id" binding:"required,maxb=255"`
	SessionID    string `json:"session_id" binding:"required,len=26"`
	ClientID     string `json:"client_id" binding:"required,maxb=255"`
	ConnectionID string `json:"connection_id" binding:"required,len=26"`

	Multistream *bool `json:"multistream" binding:"required"`
	Simulcast   *bool `json:"simulcast" binding:"required"`
	Spotlight   *bool `json:"spotlight" binding:"required"`

	Stats []json.RawMessage `json:"stats" binding:"required"`
}

func MaximumNumberOfBytesFunc(fl validator.FieldLevel) bool {
	param := fl.Param()

	// 255 バイトまで指定可能
	length, err := strconv.ParseUint(param, 10, 8)
	if err != nil {
		panic(err)
	}

	return uint64(fl.Field().Len()) <= length
}
