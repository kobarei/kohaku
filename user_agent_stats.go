package kohaku

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	db "github.com/shiguredo/kohaku/db/sqlc"
)

// TODO(v): sqlc したいが厳しそう
func (s *Server) collectorUserAgentStats(c *gin.Context, stats SoraConnectionStats) error {
	if err := s.InsertSoraConnections(c, stats); err != nil {
		return err
	}

	rtc := &RTC{
		Time:         &stats.Timestamp,
		ConnectionID: stats.ConnectionID,
	}

	for _, v := range stats.Stats {
		rtcStats := new(RTCStats)
		if err := json.Unmarshal(v, &rtcStats); err != nil {
			return err
		}

		// Type が送られてこない場合を考慮してる
		switch *rtcStats.Type {
		case "codec":
			stats := new(RTCCodecStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}

			ds := goqu.Insert("rtc_codec_stats").Rows(
				RTCCodec{
					RTC:           *rtc,
					RTCCodecStats: *stats,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := s.pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "inbound-rtp":
			stats := new(RTCInboundRTPStreamStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}

			if stats.PerDscpPacketsReceived != nil {
				// record は一旦文字列として扱う
				perDscpPacketsReceived, err := json.Marshal(stats.PerDscpPacketsReceived)
				if err != nil {
					return err
				}
				stats.PerDscpPacketsReceived = string(perDscpPacketsReceived)
			}

			ds := goqu.Insert("rtc_inbound_rtp_stream_stats").Rows(
				RTCInboundRTPStream{
					RTC:                      *rtc,
					RTCInboundRTPStreamStats: *stats,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := s.pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "outbound-rtp":
			stats := new(RTCOutboundRTPStreamStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}

			// record は一旦文字列として扱う
			if *stats.Kind == "video" {
				qualityLimitationDurations, err := json.Marshal(stats.QualityLimitationDurations)
				if err != nil {
					return err
				}
				stats.QualityLimitationDurations = string(qualityLimitationDurations)

				if stats.PerDscpPacketsSent != nil {
					perDscpPacketsSent, err := json.Marshal(stats.PerDscpPacketsSent)
					if err != nil {
						return err
					}
					stats.PerDscpPacketsSent = string(perDscpPacketsSent)
				}
			}

			ds := goqu.Insert("rtc_outbound_rtp_stream_stats").Rows(
				RTCOutboundRTPStream{
					RTC:                       *rtc,
					RTCOutboundRTPStreamStats: *stats,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := s.pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "remote-inbound-rtp":
			stats := new(RTCRemoteInboundRTPStreamStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_remote_inbound_rtp_stream_stats").Rows(
				RTCRemoteInboundRTPStream{
					RTC:                            *rtc,
					RTCRemoteInboundRTPStreamStats: *stats,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := s.pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "remote-outbound-rtp":
			stats := new(RTCRemoteOutboundRTPStreamStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_remote_outbound_rtp_stream_stats").Rows(
				RTCRemoteOutboundRTPStream{
					RTC:                             *rtc,
					RTCRemoteOutboundRTPStreamStats: *stats,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := s.pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "media-source":
			// RTCAudioSourceStats or RTCVideoSourceStats depending on its kind.
			stats := new(RTCMediaSourceStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
			switch *stats.Kind {
			case "audio":
				stats := new(RTCAudioSourceStats)
				if err := json.Unmarshal(v, &stats); err != nil {
					return err
				}
				ds := goqu.Insert("rtc_audio_source_stats").Rows(
					RTCAuidoSource{
						RTC:                 *rtc,
						RTCAudioSourceStats: *stats,
					},
				)
				insertSQL, _, _ := ds.ToSQL()
				_, err := s.pool.Exec(context.Background(), insertSQL)
				if err != nil {
					return err
				}
			case "video":
				stats := new(RTCVideoSourceStats)
				if err := json.Unmarshal(v, &stats); err != nil {
					return err
				}
				ds := goqu.Insert("rtc_video_source_stats").Rows(
					RTCVideoSource{
						RTC:                 *rtc,
						RTCVideoSourceStats: *stats,
					},
				)
				insertSQL, _, _ := ds.ToSQL()
				_, err := s.pool.Exec(context.Background(), insertSQL)
				if err != nil {
					return err
				}
			}
		case "csrc":
			stats := new(RTCRTPContributingSourceStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
		case "peer-connection":
			stats := new(RTCPeerConnectionStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
		case "data-channel":
			stats := new(RTCDataChannelStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_data_channel_stats").Rows(
				RTCDataChannel{
					RTC:                 *rtc,
					RTCDataChannelStats: *stats,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := s.pool.Exec(context.Background(), insertSQL)
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
			stats := new(RTCRTPTransceiverStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
		case "sender":
			// TODO(v): データベース書き込み
			stats := new(RTCMediaHandlerStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
			switch *stats.Kind {
			case "audio":
				stats := new(RTCAudioSenderStats)
				if err := json.Unmarshal(v, &stats); err != nil {
					return err
				}
			case "video":
				stats := new(RTCVideoSenderStats)
				if err := json.Unmarshal(v, &stats); err != nil {
					return err
				}
			}
		case "receiver":
			// TODO(v): データベース書き込み
			stats := new(RTCMediaHandlerStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
			switch *stats.Kind {
			case "audio":
				stats := new(RTCAudioReceiverStats)
				if err := json.Unmarshal(v, &stats); err != nil {
					return err
				}
			case "video":
				stats := new(RTCVideoReceiverStats)
				if err := json.Unmarshal(v, &stats); err != nil {
					return err
				}
			}
		case "transport":
			stats := new(RTCTransportStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_transport_stats").Rows(
				RTCTransport{
					RTC:               *rtc,
					RTCTransportStats: *stats,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := s.pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "sctp-transport":
			stats := new(RTCSctpTransportStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
		case "candidate-pair":
			stats := new(RTCIceCandidatePairStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_ice_candidate_pair_stats").Rows(
				RTCIceCandidatePair{
					RTC:                      *rtc,
					RTCIceCandidatePairStats: *stats,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := s.pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "local-candidate", "remote-candidate":
			stats := new(RTCIceCandidateStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
			ds := goqu.Insert("rtc_ice_candidate_stats").Rows(
				RTCIceCandidate{
					RTC:                  *rtc,
					RTCIceCandidateStats: *stats,
				},
			)
			insertSQL, _, _ := ds.ToSQL()
			_, err := s.pool.Exec(context.Background(), insertSQL)
			if err != nil {
				return err
			}
		case "certificate":
			stats := new(RTCCertificateStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
		case "ice-server":
			stats := new(RTCIceServerStats)
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
		default:
			// TODO: return err にする
			fmt.Println(rtcStats.ID)
		}

	}
	return nil
}

func (s *Server) InsertSoraConnections(ctx context.Context, stats SoraConnectionStats) error {
	if err := s.query.InsertSoraConnection(ctx, db.InsertSoraConnectionParams{
		Timestamp:    stats.Timestamp,
		Label:        stats.Label,
		Version:      stats.Version,
		NodeName:     stats.NodeName,
		Multistream:  *stats.Multistream,
		Simulcast:    *stats.Simulcast,
		Spotlight:    *stats.Spotlight,
		Role:         stats.Role,
		ChannelID:    stats.ChannelID,
		SessionID:    stats.SessionID,
		ClientID:     stats.ClientID,
		ConnectionID: stats.ConnectionID,
	}); err != nil {
		return err
	}
	return nil
}
