package revai

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"

	"github.com/gorilla/websocket"
)

// Error codes
const (
	CloseUnauthorized        = 4001
	CloseBadRequest          = 4002
	CloseInsufficientCredits = 4003
	CloseServerShuttingDown  = 4010
	CloseNoInstanceAvailable = 4013
	CloseTooManyRequests     = 4029
)

// Whether or not connection should be retried
var shouldErrorRetry = map[int]bool{
	CloseUnauthorized:        false,
	CloseBadRequest:          false,
	CloseInsufficientCredits: false,
	CloseServerShuttingDown:  true,
	CloseNoInstanceAvailable: true,
	CloseTooManyRequests:     false,
}

// Whether or not connection should be retried
var shouldErrorRetry = map[int]bool{
	CloseUnauthorized:        "Unauthorized. The provided access token is invalid.",
	CloseBadRequest:          "Bad request. The connection’s content-type is invalid, metadata contains too many characters or the custom vocabulary does not exist with that id.",
	CloseInsufficientCredits: "Insufficient credits. The client does not have enough credits to continue the streaming session.",
	CloseServerShuttingDown:  "Server shutting down. The connection was terminated due to the server shutting down.",
	CloseNoInstanceAvailable: "No instance available. No available streaming instances were found. User should attempt to retry the connection later.",
	CloseTooManyRequests:     "Too many requests. The number of concurrent connections exceeded the limit. Contact customer support to increase it.",
}

// A close message from rev see https://www.rev.ai/docs/streaming#section/Error-Codes
type StreamingError struct {
	// Error code
	Code int

	// The error string
	Text string

	// Whether or not the connection should be retried
	ShouldRetry bool
}

func (e *StreamingError) Error() string {
	return fmt.Sprintf("Streaming error: %s", e.Text)
}

// StreamService provides access to the stream related functions
// in the Rev.ai API.
type StreamService service

// StreamMessage represents a rev.ai websocket stream message.
type StreamMessage struct {
	Type     string    `json:"type"`
	Ts       float64   `json:"ts"`
	EndTs    float64   `json:"end_ts"`
	Elements []Element `json:"elements"`
}

// Conn represents a websocket connection to the Rev.ai Api.
// It has certain helper methods to easily parse and communicate to the
// web socket connection
type Conn struct {
	Msg chan StreamMessage

	conn *websocket.Conn
}

// Write sends a message to the websocket connection
func (c *Conn) Write(r io.Reader) error {
	w, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return fmt.Errorf("failed getting writer %w", err)
	}

	if _, err := io.Copy(w, r); err != nil {
		return fmt.Errorf("failed to copy data %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to send data")
	}

	return nil
}

// Send EOS to let Rev know we are done. see https://www.rev.ai/docs/streaming#section/Client-to-Rev.ai-Input/Sending-Audio-to-Rev.ai
func (c *Conn) WriteDone() error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte("EOS"))
}

// Close closes the message chan and the websocket connection
func (c *Conn) Close() error {
	close(c.Msg)

	return c.conn.Close()
}

// DialStreamParams specifies the parameters to the
// StreamService.Dial method.
type DialStreamParams struct {
	ContentType        string
	Metadata           string
	FilterProfanity    bool
	RemoveDisfluencies bool
	CustomVocabularyID string
}

type dialStreamParams struct {
	ContentType        string `url:"content_type"`
	Metadata           string `url:"metadata,omitempty"`
	RemoveDisfluencies bool   `url:"remove_disfluencies,omitempty"`
	FilterProfanity    bool   `url:"filter_profanity"`
	CustomVocabularyID string `url:"custom_vocabulary_id"`
	AccessToken        string `url:"access_token"`
}

// Dial dials a WebSocket request to the Rev.ai Streaming api.
// https://www.rev.ai/docs/streaming#section/Overview
func (s *StreamService) Dial(ctx context.Context, params *DialStreamParams) (*Conn, error) {
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}

	u, err := s.streamURL(params)
	if err != nil {
		return nil, fmt.Errorf("failed creating url %w", err)
	}

	websocketConn, _, err := dialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed dialing %w", err)
	}

	conn := &Conn{
		conn: websocketConn,
		Msg:  make(chan StreamMessage),
	}

	go func() {
		defer conn.Close()
		for {
			var msg StreamMessage
			if err := conn.conn.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					return
				}
				continue
			}
			conn.Msg <- msg
		}
	}()

	return conn, nil
}

func (s *StreamService) streamURL(params *DialStreamParams) (*url.URL, error) {
	rel := &url.URL{Scheme: "wss", Path: "/speechtotext/v1/stream", Host: s.client.BaseURL.Host}

	p := &dialStreamParams{
		AccessToken:        s.client.APIKey,
		ContentType:        params.ContentType,
		Metadata:           params.Metadata,
		FilterProfanity:    params.FilterProfanity,
		RemoveDisfluencies: params.RemoveDisfluencies,
		CustomVocabularyID: params.CustomVocabularyID,
	}

	v, err := query.Values(p)
	if err != nil {
		return nil, err
	}

	rel.RawQuery = v.Encode()

	return rel, nil
}
