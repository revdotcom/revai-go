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
	ErrCloseUnauthorized        = 4001
	ErrCloseBadRequest          = 4002
	ErrCloseInsufficientCredits = 4003
	ErrCloseServerShuttingDown  = 4010
	ErrCloseNoInstanceAvailable = 4013
	ErrCloseTooManyRequests     = 4029
)

// Whether or not connection should be retried
var shouldErrorRetry = map[int]bool{
	ErrCloseUnauthorized:        false,
	ErrCloseBadRequest:          false,
	ErrCloseInsufficientCredits: false,
	ErrCloseServerShuttingDown:  true,
	ErrCloseNoInstanceAvailable: true,
	ErrCloseTooManyRequests:     false,
}

// Whether or not connection should be retried
var errorMsgs = map[int]string{
	ErrCloseUnauthorized:        "Unauthorized. The provided access token is invalid.",
	ErrCloseBadRequest:          "Bad request. The connection’s content-type is invalid, metadata contains too many characters or the custom vocabulary does not exist with that id.",
	ErrCloseInsufficientCredits: "Insufficient credits. The client does not have enough credits to continue the streaming session.",
	ErrCloseServerShuttingDown:  "Server shutting down. The connection was terminated due to the server shutting down.",
	ErrCloseNoInstanceAvailable: "No instance available. No available streaming instances were found. User should attempt to retry the connection later.",
	ErrCloseTooManyRequests:     "Too many requests. The number of concurrent connections exceeded the limit. Contact customer support to increase it.",
}

// A close message from rev see https://www.rev.ai/docs/streaming#section/Error-Codes
type RevError struct {
	// Error code
	Code int

	// The error string
	Text string
}

// A retriable eror
type RetriableError struct {
	// Error code
	Code int

	// The error string
	Text string
}

func (e RevError) Error() string {
	return fmt.Sprintf("Streaming error: %s", e.Text)
}

func (e RetriableError) Error() string {
	return fmt.Sprintf("Retriable streaming error: %s", e.Text)
}

// Check if the code is a Rev error if so return it.
func IsRevError(code int) (bool, error) {
	errorString, exists := errorMsgs[code]
	if exists {
		shouldRetry := shouldErrorRetry[code]
		if shouldRetry {
			return true, RetriableError{code, errorString}
		} else {
			return true, RevError{code, errorString}
		}
	} else {
		return false, nil
	}
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
	msg  chan StreamMessage
	err  chan error
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

// Recv get messages back from rev
func (c *Conn) Recv() (StreamMessage, error) {
	select {
	case err := <-c.err:
		return nil, err
	case msg := <-c.msg:
		return msg, nil
	}
}

// Send EOS to let Rev know we are done. see https://www.rev.ai/docs/streaming#section/Client-to-Rev.ai-Input/Sending-Audio-to-Rev.ai
func (c *Conn) WriteDone() error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte("EOS"))
}

// Close closes the message chan and the websocket connection
func (c *Conn) Close() error {
	close(c.msg)

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
		msg:  make(chan StreamMessage),
		err:  make(chan error),
	}

	go func() {
		defer conn.Close()
		for {
			var msg StreamMessage
			if err := conn.conn.ReadJSON(&msg); err != nil {
				if e, ok := err.(*websocket.CloseError); ok {
					if isRevError, revError := IsRevError(e.Code); isRevError {
						conn.err <- revError
						return
					}
				}

				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					// perhaps an error should be sent on Err here too
					conn.err <- err
					return
				}
				continue
			}
			conn.msg <- msg
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
