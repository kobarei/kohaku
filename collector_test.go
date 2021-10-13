package kohaku

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

// TODO(v): mockDB を導入する

var (
	collectorTypeCodecJSON = `{
    "role": "sendrecv",
    "type": "connection.remote",
    "channel_id": "sora",
    "client_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "connection_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "session_id": "KE9C2QKV892TD03CA2CR38BV4G",
    "stats": [{
      "id": "RTCCodec_video_V04mIx_Inbound_120",
      "timestamp": 1628869622194.298,
      "type": "codec",
      "transportId": "RTCTransport_data_1",
      "payloadType": 120,
      "mimeType": "video/VP9",
      "clockRate": 90000,
      "sdpFmtpLine": "profile-id=0"
    }],
    "multistream": false,
    "spotlight": false,
    "simulcast": false,
    "timestamp":"2021-09-24T08:15:31.854427Z",
    "version":"2021.2-canary.23"}
  }`

	collectorTypeOutboundRTPJSON = `{
    "role": "sendrecv",
    "type": "connection.remote",
    "channel_id": "sora",
    "client_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "connection_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "session_id": "KE9C2QKV892TD03CA2CR38BV4G",
    "stats": [{
      "id": "RTCOutboundRTPVideoStream_1028062523",
      "timestamp": 1628927446077.817,
      "type": "outbound-rtp",
      "ssrc": 1028062523,
      "kind": "video",
      "trackId": "RTCMediaStreamTrack_sender_2",
      "transportId": "RTCTransport_data_1",
      "codecId": "RTCCodec_video_oVLkJT_Outbound_120",
      "mediaType": "video",
      "mediaSourceId": "RTCVideoSource_2",
      "remoteId": "RTCRemoteInboundRtpVideoStream_1028062523",
      "packetsSent": 2056,
      "retransmittedPacketsSent": 0,
      "bytesSent": 2059458,
      "headerBytesSent": 54692,
      "retransmittedBytesSent": 0,
      "framesEncoded": 538,
      "keyFramesEncoded": 1,
      "totalEncodeTime": 2.541,
      "totalEncodedBytesTarget": 0,
      "frameWidth": 640,
      "frameHeight": 480,
      "framesPerSecond": 30,
      "framesSent": 538,
      "hugeFramesSent": 0,
      "totalPacketSendDelay": 60.249,
      "qualityLimitationReason": "none",
      "qualityLimitationResolutionChanges": 0,
      "encoderImplementation": "libvpx",
      "firCount": 0,
      "pliCount": 0,
      "nackCount": 0,
      "qpSum": 40927
    }],
    "multistream": false,
    "spotlight": false,
    "simulcast": false,
    "timestamp":"2021-09-24T08:15:31.854427Z",
    "version":"2021.2-canary.23"}
  }`

	collectorTypeMediaSourceJSON = `{
    "role": "sendrecv",
    "channel_id":"sora",
    "client_id":"KB0DR2FWT13C70S0NYS11P04C0",
    "connection_id":"KB0DR2FWT13C70S0NYS11P04C0",
    "session_id": "KE9C2QKV892TD03CA2CR38BV4G",
    "id":"3Q1Y9Y6B9X7CKDXFWNZX3PVJ9W",
    "label":"WebRTC.SFU.Sora",
    "stats":[{
      "id":"RTCAudioSource_3",
      "kind":"audio",
      "timestamp":1.63247133184561206055e+12,
      "type":"media-source",
      "audioLevel":2.13629566331980345895e-03,
      "echoReturnLoss":-30,
      "echoReturnLossEnhancement":1.75512030720710754395e-01,
      "totalAudioEnergy":4.43774538826569622503e-05,
      "totalSamplesDuration":1.00599999999998299671e+01,
      "trackIdentifier":"e6763fb2-0f7a-46f2-af0d-bdde1fcc1ee7"
    },{
      "height":480,
      "id":"RTCVideoSource_4",
      "kind":"video",
      "timestamp":1.63247133184561206055e+12,
      "type":"media-source",
      "width":640,
      "frames":284,
      "framesPerSecond":30,
      "trackIdentifier":"9d2bc9fd-361e-4a74-9030-ee2212aadfee"
    }],
    "multistream": false,
    "spotlight": false,
    "simulcast": false,
    "timestamp":"2021-09-24T08:15:31.854427Z",
    "type":"connection.remote",
    "version":"2021.2-canary.23"
  }
`

	collectorTypeDataChannelJSON = `{
    "role": "sendrecv",
    "channel_id":"sora",
    "client_id":"KB0DR2FWT13C70S0NYS11P04C0",
    "connection_id":"KB0DR2FWT13C70S0NYS11P04C0",
    "session_id": "KE9C2QKV892TD03CA2CR38BV4G",
    "id":"3Q1Y9Y6B9X7CKDXFWNZX3PVJ9W",
    "label":"WebRTC.SFU.Sora",
    "stats":[{
    },{
      "id":"RTCDataChannel_10",
      "label":"stats",
      "protocol":"",
      "state":"open",
      "timestamp":1.63247133184561206055e+12,
      "type":"data-channel",
      "bytesReceived":56,
      "bytesSent":2853,
      "dataChannelIdentifier":8,
      "messagesReceived":2,
      "messagesSent":1
    },{
      "id":"RTCDataChannel_6",
      "label":"signaling",
      "protocol":"",
      "state":"open",
      "timestamp":1.63247133184561206055e+12,
      "type":"data-channel",
      "bytesReceived":0,
      "bytesSent":0,
      "dataChannelIdentifier":0,
      "messagesReceived":0,
      "messagesSent":0
    },{
      "id":"RTCDataChannel_7",
      "label":"notify",
      "protocol":"",
      "state":"open",
      "timestamp":1.63247133184561206055e+12,
      "type":"data-channel",
      "bytesReceived":64,
      "bytesSent":0,
      "dataChannelIdentifier":2,
      "messagesReceived":1,
      "messagesSent":0
    },{
      "id":"RTCDataChannel_8",
      "label":"push",
      "protocol":"",
      "state":"open",
      "timestamp":1.63247133184561206055e+12,
      "type":"data-channel",
      "bytesReceived":0,
      "bytesSent":0,
      "dataChannelIdentifier":4,
      "messagesReceived":0,
      "messagesSent":0
    },{
      "id":"RTCDataChannel_9",
      "label":"e2ee",
      "protocol":"",
      "state":"open",
      "timestamp":1.63247133184561206055e+12,
      "type":"data-channel",
      "bytesReceived":0,
      "bytesSent":0,
      "dataChannelIdentifier":6,
      "messagesReceived":0,
      "messagesSent":0
    }],
    "multistream": false,
    "spotlight": false,
    "simulcast": false,
    "timestamp":"2021-09-24T08:15:31.854427Z",
    "type":"connection.remote",
    "version":"2021.2-canary.23"
  }
`

	collectorTypeCandidatePairJSON = `{
    "role": "sendrecv",
    "channel_id":"sora",
    "client_id":"KB0DR2FWT13C70S0NYS11P04C0",
    "connection_id":"KB0DR2FWT13C70S0NYS11P04C0",
    "session_id": "KE9C2QKV892TD03CA2CR38BV4G",
    "id":"3Q1Y9Y6B9X7CKDXFWNZX3PVJ9W",
    "label":"WebRTC.SFU.Sora",
    "stats":[{
      "id":"RTCIceCandidatePair_zNnR\/QQb_SwzcXtlY",
      "priority":179616219446525440,
      "state":"succeeded",
      "timestamp":1.63247133184561206055e+12,
      "type":"candidate-pair",
      "writable":true,
      "availableOutgoingBitrate":500000,
      "bytesReceived":3437,
      "bytesSent":670907,
      "consentRequestsSent":7,
      "currentRoundTripTime":1.00000000000000002082e-03,
      "localCandidateId":"RTCIceCandidate_zNnR\/QQb",
      "nominated":true,
      "remoteCandidateId":"RTCIceCandidate_SwzcXtlY",
      "requestsReceived":9,
      "requestsSent":1,
      "responsesReceived":8,
      "responsesSent":9,
      "totalRoundTripTime":5.00000000000000010408e-03,
      "transportId":"RTCTransport_data_1"
    }],
    "multistream": false,
    "spotlight": false,
    "simulcast": false,
    "timestamp":"2021-09-24T08:15:31.854427Z",
    "type":"connection.remote",
    "version":"2021.2-canary.23"
  }
`

	collectorTypeRemoteInboundRTPJSON = `{
    "role": "sendrecv",
    "channel_id":"sora",
    "client_id":"KB0DR2FWT13C70S0NYS11P04C0",
    "connection_id":"KB0DR2FWT13C70S0NYS11P04C0",
    "session_id": "KE9C2QKV892TD03CA2CR38BV4G",
    "id":"3Q1Y9Y6B9X7CKDXFWNZX3PVJ9W",
    "label":"WebRTC.SFU.Sora",
    "stats":[{
      "fractionLost":0,
      "id":"RTCRemoteInboundRtpAudioStream_1073653878",
      "kind":"audio",
      "ssrc":1073653878,
      "timestamp":1.63247133183873803711e+12,
      "type":"remote-inbound-rtp",
      "codecId":"RTCCodec_audio_Z4hWio_Outbound_109",
      "jitter":2.91666666666666690808e-04,
      "localId":"RTCOutboundRTPAudioStream_1073653878",
      "packetsLost":0,
      "roundTripTime":1.00000000000000002082e-03,
      "roundTripTimeMeasurements":2,
      "totalRoundTripTime":2.00000000000000004163e-03,
      "transportId":"RTCTransport_data_1"
    },{
      "fractionLost":0,
      "id":"RTCRemoteInboundRtpVideoStream_1597059872",
      "kind":"video",
      "ssrc":1597059872,
      "timestamp":1.63247133144896899414e+12,
      "type":"remote-inbound-rtp",
      "codecId":"RTCCodec_video_Xn86YA_Outbound_120",
      "jitter":2.34444444444444500403e-03,
      "localId":"RTCOutboundRTPVideoStream_1597059872",
      "packetsLost":0,
      "roundTripTime":1.00000000000000002082e-03,
      "roundTripTimeMeasurements":10,
      "totalRoundTripTime":1.00000000000000002082e-02,
      "transportId":"RTCTransport_data_1"
    }],
    "multistream": false,
    "spotlight": false,
    "simulcast": false,
    "timestamp":"2021-09-24T08:15:31.854427Z",
    "type":"connection.remote",
    "version":"2021.2-canary.23"
  }
`

	collectorTypeTransportJSON = `{
    "role": "sendrecv",
    "channel_id":"sora",
    "client_id":"KB0DR2FWT13C70S0NYS11P04C0",
    "connection_id":"KB0DR2FWT13C70S0NYS11P04C0",
    "session_id": "KE9C2QKV892TD03CA2CR38BV4G",
    "id":"3Q1Y9Y6B9X7CKDXFWNZX3PVJ9W",
    "label":"WebRTC.SFU.Sora",
    "stats":[{
      "id":"RTCTransport_data_1",
      "timestamp":1.63247133184561206055e+12,
      "type":"transport",
      "bytesReceived":3437,
      "bytesSent":670907,
      "dtlsCipher":"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
      "dtlsState":"connected",
      "localCertificateId":"RTCCertificate_60:96:2D:B7:B6:D8:0A:15:92:45:21:4E:1B:DB:66:01:CC:44:65:D2:43:44:31:15:E4:09:D0:64:58:A2:BF:84",
      "packetsReceived":34,
      "packetsSent":1105,
      "remoteCertificateId":"RTCCertificate_1C:64:46:83:99:7F:9C:44:8A:5B:5C:DF:07:C0:A4:2D:39:51:72:39:8B:76:B5:D2:75:0C:A4:0D:58:FF:67:69",
      "selectedCandidatePairChanges":1,
      "selectedCandidatePairId":"RTCIceCandidatePair_zNnR\/QQb_SwzcXtlY",
      "srtpCipher":"AEAD_AES_128_GCM",
      "tlsVersion":"FEFD"
    }],
    "multistream": false,
    "spotlight": false,
    "simulcast": false,
    "timestamp":"2021-09-24T08:15:31.854427Z",
    "type":"connection.remote",
    "version":"2021.2-canary.23"
  }
`
)

var (
	invalidConnectionIDLengthJSON = `{
    "role": "sendrecv",
    "type": "connection.remote",
    "channel_id": "sora",
    "client_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "connection_id": "2QB23E50YD6FKEFG9GW2TX86RC===",
    "session_id": "KE9C2QKV892TD03CA2CR38BV4G",
    "stats": [{
      "id": "RTCCodec_video_V04mIx_Inbound_120",
      "timestamp": 1628869622194.298,
      "type": "codec",
      "transportId": "RTCTransport_data_1",
      "payloadType": 120,
      "mimeType": "video/VP9",
      "clockRate": 90000,
      "sdpFmtpLine": "profile-id=0"
    }],
    "multistream": false,
    "spotlight": false,
    "simulcast": false,
    "timestamp":"2021-09-24T08:15:31.854427Z",
    "version":"2021.2-canary.23"}
  }`

	unexpectedTypeJSON = `{
    "role": "sendrecv",
    "type": "connection.unexpected_type",
    "channel_id": "sora",
    "client_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "connection_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "session_id": "KE9C2QKV892TD03CA2CR38BV4G",
    "stats": [{
      "id": "RTCCodec_video_V04mIx_Inbound_120",
      "timestamp": 1628869622194.298,
      "type": "codec",
      "transportId": "RTCTransport_data_1",
      "payloadType": 120,
      "mimeType": "video/VP9",
      "clockRate": 90000,
      "sdpFmtpLine": "profile-id=0"
    }],
    "multistream": false,
    "spotlight": false,
    "simulcast": false,
    "timestamp":"2021-09-24T08:15:31.854427Z",
    "version":"2021.2-canary.23"}
  }`

	missingTimestampJSON = `{
    "role": "sendrecv",
    "type": "connection.unexpected_type",
    "channel_id": "sora",
    "client_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "connection_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "session_id": "KE9C2QKV892TD03CA2CR38BV4G",
    "stats": [{
      "id": "RTCCodec_video_V04mIx_Inbound_120",
      "timestamp": 1628869622194.298,
      "type": "codec",
      "transportId": "RTCTransport_data_1",
      "payloadType": 120,
      "mimeType": "video/VP9",
      "clockRate": 90000,
      "sdpFmtpLine": "profile-id=0"
    }],
    "multistream": false,
    "spotlight": false,
    "simulcast": false,
    "version":"2021.2-canary.23"}
  }`

	invalidChannelIDLengthJSON = `{
    "role": "sendrecv",
    "type": "connection.remote",
    "channel_id": "2QB23E50YD6FKEFG9GW2TX86RC2QB23E50YD6FKEFG9GW2TX86RC2QB23E50YD6FKEFG9GW2TX86RC2QB23E50YD6FKEFG9GW2TX86RC2QB23E50YD6FKEFG9GW2TX86RC2QB23E50YD6FKEFG9GW2TX86RC2QB23E50YD6FKEFG9GW2TX86RC2QB23E50YD6FKEFG9GW2TX86RC2QB23E50YD6FKEFG9GW2TX86RC2QB23E50YD6FKEFG9GW2TX",
    "client_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "connection_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "session_id": "KE9C2QKV892TD03CA2CR38BV4G",
    "stats": [{
      "id": "RTCCodec_video_V04mIx_Inbound_120",
      "timestamp": 1628869622194.298,
      "type": "codec",
      "transportId": "RTCTransport_data_1",
      "payloadType": 120,
      "mimeType": "video/VP9",
      "clockRate": 90000,
      "sdpFmtpLine": "profile-id=0"
    }],
    "multistream": false,
    "spotlight": false,
    "simulcast": false,
    "timestamp":"2021-09-24T08:15:31.854427Z",
    "version":"2021.2-canary.23"}
  }`
)

const (
	connStr     = "postgres://postgres:password@127.0.0.1:5432/%s?sslmode=disable"
	dbName      = "kohakutest"
	sqlFilePath = "scripts/timescaledb.sql"

	channelID    = "sora"
	connectionID = "KB0DR2FWT13C70S0NYS11P04C0"
	clientID     = "KB0DR2FWT13C70S0NYS11P04C0"
)

var (
	pool               *pgxpool.Pool
	postgresDBURL      = fmt.Sprintf(connStr, "postgres")
	kohakuDBURL        = fmt.Sprintf(connStr, dbName)
	dropDBSQL          = fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	createDBSQL        = fmt.Sprintf("CREATE DATABASE %s", dbName)
	createExtensionSQL = fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS timescaledb")
	server             *Server
)

func createTable() error {
	return exec.Command("psql", "-d", dbName, "-f", sqlFilePath).Run()
}

func getStatsType(table, connectionID string) (*string, error) {
	selectSQL := fmt.Sprintf("SELECT stats_type FROM %s WHERE sora_connection_id=$1", table)
	row := pool.QueryRow(context.Background(), selectSQL, connectionID)

	var statsType string
	if err := row.Scan(&statsType); err != nil {
		return nil, err
	}

	return &statsType, nil
}

func TestMain(m *testing.M) {
	// DB の削除、作成、Table の作成
	postgresDBConfig, err := pgxpool.ParseConfig(postgresDBURL)
	if err != nil {
		panic(err)
	}
	postgresDBPool, err := pgxpool.ConnectConfig(context.Background(), postgresDBConfig)
	if err != nil {
		panic(err)
	}
	defer postgresDBPool.Close()

	_, err = postgresDBPool.Exec(context.Background(), dropDBSQL)
	if err != nil {
		panic(err)
	}
	_, err = postgresDBPool.Exec(context.Background(), createDBSQL)
	if err != nil {
		panic(err)
	}
	config, err := pgxpool.ParseConfig(kohakuDBURL)
	if err != nil {
		panic(err)
	}
	pool, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}

	_, err = pool.Exec(context.Background(), createExtensionSQL)
	if err != nil {
		panic(err)
	}

	// TODO: pool を使用する
	if err := createTable(); err != nil {
		panic(err)
	}

	server = NewServer(&KohakuConfig{}, pool)

	status := m.Run()

	// DB の削除
	pool.Close()
	_, err = postgresDBPool.Exec(context.Background(), dropDBSQL)
	if err != nil {
		panic(err)
	}

	os.Exit(status)
}

func TestTypeOutboundRTPCollector(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(collectorTypeOutboundRTPJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	server.Collector(c)
	assert.Equal(t, http.StatusNoContent, c.Writer.Status())

	statsType, err := getStatsType("rtc_outbound_rtp_stream_stats", "2QB23E50YD6FKEFG9GW2TX86RC")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "outbound-rtp", *statsType)
}

func TestTypeCodecCollector(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(collectorTypeCodecJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	server.Collector(c)
	assert.Equal(t, http.StatusNoContent, c.Writer.Status())

	statsType, err := getStatsType("rtc_codec_stats", "2QB23E50YD6FKEFG9GW2TX86RC")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "codec", *statsType)
}

func TestTypeMediaSourceCollector(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(collectorTypeMediaSourceJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	server.Collector(c)
	assert.Equal(t, http.StatusNoContent, c.Writer.Status())

	statsType, err := getStatsType("rtc_audio_source_stats", connectionID)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "media-source", *statsType)
}

func TestTypeDataChannelCollector(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(collectorTypeDataChannelJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	server.Collector(c)
	assert.Equal(t, http.StatusNoContent, c.Writer.Status())

	statsType, err := getStatsType("rtc_data_channel_stats", connectionID)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "data-channel", *statsType)
}

func TestTypeCandidatePairCollector(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(collectorTypeCandidatePairJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	server.Collector(c)
	assert.Equal(t, http.StatusNoContent, c.Writer.Status())

	statsType, err := getStatsType("rtc_ice_candidate_pair_stats", connectionID)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "candidate-pair", *statsType)
}

func TestTypeRemoteInboundRTPCollector(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(collectorTypeRemoteInboundRTPJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	server.Collector(c)
	assert.Equal(t, http.StatusNoContent, c.Writer.Status())

	statsType, err := getStatsType("rtc_remote_inbound_rtp_stream_stats", connectionID)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "remote-inbound-rtp", *statsType)
}

func TestTypeTransportCollector(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(collectorTypeTransportJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	server.Collector(c)
	assert.Equal(t, http.StatusNoContent, c.Writer.Status())

	statsType, err := getStatsType("rtc_transport_stats", connectionID)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "transport", *statsType)
}

func TestInvalidConnectionIDLength(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(invalidConnectionIDLengthJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	server.Collector(c)
	resp := rec.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	assert.NotEmpty(t, body)
}

func TestUnexpectedType(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(unexpectedTypeJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	server.Collector(c)
	resp := rec.Result()

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	assert.Empty(t, body)
}

func TestMissingTimestamp(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(missingTimestampJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	server.Collector(c)
	resp := rec.Result()

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, `{"error":"Key: 'SoraStatsExporter.Timestamp' Error:Field validation for 'Timestamp' failed on the 'required' tag"}`, string(body))
}

func TestInvalidChannelIDLength(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(invalidChannelIDLengthJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	server.Collector(c)
	resp := rec.Result()

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, `{"error":"Key: 'SoraStatsExporter.ChannelID' Error:Field validation for 'ChannelID' failed on the 'maxb' tag"}`, string(body))
}
