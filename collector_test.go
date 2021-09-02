package kohaku

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TODO(v): mockDB を導入する

var (
	collectorTypeCodecJSON = `{
    "type": "connection.peer",
    "channel_id": "sora",
    "client_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "connection_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "stats": [{
      "id": "RTCCodec_video_V04mIx_Inbound_120",
      "timestamp": 1628869622194.298,
      "type": "codec",
      "transportId": "RTCTransport_data_1",
      "payloadType": 120,
      "mimeType": "video/VP9",
      "clockRate": 90000,
      "sdpFmtpLine": "profile-id=0"
    }]
  }`

	collectorTypeOutboundRTPJSON = `{
    "type": "connection.peer",
    "channel_id": "sora",
    "client_id": "2QB23E50YD6FKEFG9GW2TX86RC",
    "connection_id": "2QB23E50YD6FKEFG9GW2TX86RC",
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
    }]
  }`
)

func TestTypeOutboundRTPCollector(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(collectorTypeCodecJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	Collector(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestTypeCodecCollector(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/collector", strings.NewReader(collectorTypeOutboundRTPJSON))
	req.Header.Set("content-type", "application/json")
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Assertions
	Collector(c)
	assert.Equal(t, http.StatusOK, rec.Code)
}
