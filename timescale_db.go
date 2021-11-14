package kohaku

import "time"

// FIXME(v): 名前はすべて仮です
type SoraNode struct {
	Timestamp time.Time `db:"timestamp"`

	Label    string `db:"label"`
	Version  string `db:"version"`
	NodeName string `db:"node_name"`
}

type SoraConnection struct {
	SoraNode

	Multistream bool `db:"multistream"`
	Simulcast   bool `db:"simulcast"`
	Spotlight   bool `db:"spotlight"`

	Role         string `db:"role"`
	ChannelID    string `db:"channel_id"`
	SessionID    string `db:"session_id"`
	ClientID     string `db:"client_id"`
	ConnectionID string `db:"connection_id"`
}

type SoraNodeErlangVM struct {
	SoraNode

	TotalMemory        uint64 `db:"total_memory"`
	TotalProcesses     uint64 `db:"total_processes"`
	TotalProcessesUsed uint64 `db:"total_processes_used"`
	TotalSystem        uint64 `db:"total_system"`
	TotalAtom          uint64 `db:"total_atom"`
	TotalAtomUsed      uint64 `db:"total_atom_used"`
	TotalBinary        uint64 `db:"total_binary"`
	TotalCode          uint64 `db:"total_code"`
	TotalETS           uint64 `db:"total_ets"`
}

type RTC struct {
	Time *time.Time `db:"time"`

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
