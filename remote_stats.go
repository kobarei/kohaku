package kohaku

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TODO(v): リファクタリング
func CollectorRemoteStats(pool *pgxpool.Pool, exporter SoraStatsExporter) error {
	rtc := &RTC{
		Time:         exporter.Timestamp,
		ChannelID:    exporter.ChannelID,
		ClientID:     exporter.ClientID,
		ConnectionID: exporter.ConnectionID,

		Label:   exporter.Label,
		Version: exporter.Version,
	}

	for _, v := range exporter.Stats {
		stats := new(RTCStats)
		if err := json.Unmarshal(v, &stats); err != nil {
			return err
		}

		// その後 type をみて struct をさらに別途デコードする
		switch stats.Type {
		case "codec":
			s := new(RTCCodecStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}

			ds := goqu.Insert("rtc_codec_stats").Rows(
				RTCCodec{
					RTC:           *rtc,
					RTCCodecStats: *s,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "inbound-rtp":
			s := new(RTCInboundRtpStreamStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}

			ds := goqu.Insert("rtc_inbound_rtp_stream_stats").Rows(
				RTCInboundRtpStream{
					RTC:                      *rtc,
					RTCInboundRtpStreamStats: *s,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "outbound-rtp":
			s := new(RTCOutboundRtpStreamStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_outbound_rtp_stream_stats").Rows(
				RTCOutboundRtpStream{
					RTC:                       *rtc,
					RTCOutboundRtpStreamStats: *s,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "remote-inbound-rtp":
			s := new(RTCRemoteInboundRtpStreamStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_remote_inbound_rtp_stream_stats").Rows(
				RTCRemoteInboundRtpStream{
					RTC:                            *rtc,
					RTCRemoteInboundRtpStreamStats: *s,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "remote-outbound-rtp":
			s := new(RTCRemoteOutboundRtpStreamStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_remote_outbound_rtp_stream_stats").Rows(
				RTCRemoteOutboundRtpStream{
					RTC:                             *rtc,
					RTCRemoteOutboundRtpStreamStats: *s,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "media-source":
			// RTCAudioSourceStats or RTCVideoSourceStats depending on its kind.
			s := new(RTCMediaSourceStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			switch s.Kind {
			case "audio":
				s := new(RTCAudioSourceStats)
				if err := json.Unmarshal(v, &s); err != nil {
					return err
				}
				ds := goqu.Insert("rtc_audio_source_stats").Rows(
					RTCAuidoSource{
						RTC:                 *rtc,
						RTCAudioSourceStats: *s,
					},
				)
				insertSQL, _, _ := ds.ToSQL()
				_, err := pool.Exec(context.Background(), insertSQL)
				if err != nil {
					return err
				}
			case "video":
				s := new(RTCVideoSourceStats)
				if err := json.Unmarshal(v, &s); err != nil {
					return err
				}
				ds := goqu.Insert("rtc_video_source_stats").Rows(
					RTCVideoSource{
						RTC:                 *rtc,
						RTCVideoSourceStats: *s,
					},
				)
				insertSQL, _, _ := ds.ToSQL()
				_, err := pool.Exec(context.Background(), insertSQL)
				if err != nil {
					return err
				}
			}
		case "csrc":
			s := new(RTCRtpContributingSourceStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
		case "peer-connection":
			s := new(RTCPeerConnectionStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
		case "data-channel":
			s := new(RTCDataChannelStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_data_channel_stats").Rows(
				RTCDataChannel{
					RTC:                 *rtc,
					RTCDataChannelStats: *s,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "stream":
			// Obsolete stats
			return nil
		case "track":
			// Obsolete stats
			return nil
		case "transceiver":
			// TODO(v): データベース書き込み
			s := new(RTCRtpTransceiverStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
		case "sender":
			// TODO(v): データベース書き込み
			s := new(RTCMediaHandlerStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			switch s.Kind {
			case "audio":
				s := new(RTCAudioSenderStats)
				if err := json.Unmarshal(v, &s); err != nil {
					return err
				}
			case "video":
				s := new(RTCVideoSenderStats)
				if err := json.Unmarshal(v, &s); err != nil {
					return err
				}
			}
		case "receiver":
			// TODO(v): データベース書き込み
			s := new(RTCMediaHandlerStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			switch s.Kind {
			case "audio":
				s := new(RTCAudioReceiverStats)
				if err := json.Unmarshal(v, &s); err != nil {
					return err
				}
			case "video":
				s := new(RTCVideoReceiverStats)
				if err := json.Unmarshal(v, &s); err != nil {
					return err
				}
			}
		case "transport":
			s := new(RTCTransportStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_transport_stats").Rows(
				RTCTransport{
					RTC:               *rtc,
					RTCTransportStats: *s,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "sctp-transport":
			s := new(RTCSctpTransportStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
		case "candidate-pair":
			s := new(RTCIceCandidatePairStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_ice_candidate_pair_stats").Rows(
				RTCIceCandidatePair{
					RTC:                      *rtc,
					RTCIceCandidatePairStats: *s,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "local-candidate", "remote-candidate":
			s := new(RTCIceCandidateStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_ice_candidate_stats").Rows(
				RTCIceCandidate{
					RTC:                  *rtc,
					RTCIceCandidateStats: *s,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "certificate":
			s := new(RTCCertificateStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
		case "ice-server":
			s := new(RTCIceServerStats)
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
		default:
			fmt.Println(stats.ID)
		}

	}
	return nil
}
