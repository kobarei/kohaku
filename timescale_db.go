package kohaku

import "time"

// FIXME(v): 名前はすべて仮です

type SoraConnections struct {
	Time time.Time `db:"time"`

	ChannelID    string `db:"sora_channel_id"`
	ClientID     string `db:"sora_client_id"`
	ConnectionID string `db:"sora_connection_id"`

	Label   string `db:"sora_label"`
	Version string `db:"sora_version"`
}

type RTC struct {
	Time time.Time `db:"time"`

	ConnectionID string `db:"sora_connection_id"`
}

type RTCCodec struct {
	RTC
	RTCCodecStats
}

type RTCInboundRtpStream struct {
	RTC
	RTCInboundRtpStreamStats
}

type RTCRemoteInboundRtpStream struct {
	RTC
	RTCRemoteInboundRtpStreamStats
}

type RTCOutboundRtpStream struct {
	RTC
	RTCOutboundRtpStreamStats
}

type RTCRemoteOutboundRtpStream struct {
	RTC
	RTCRemoteOutboundRtpStreamStats
}

type RTCAuidoSource struct {
	RTC
	RTCAudioSourceStats
}

type RTCVideoSource struct {
	RTC
	RTCVideoSourceStats
}

type RTCDataChannel struct {
	RTC
	RTCDataChannelStats
}

type RTCTransport struct {
	RTC
	RTCTransportStats
}

type RTCIceCandidate struct {
	RTC
	RTCIceCandidateStats
}

type RTCIceCandidatePair struct {
	RTC
	RTCIceCandidatePairStats
}
