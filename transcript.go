package revai

import (
	"context"
	"fmt"
	"net/http"
)

// TranscriptService provides access to the transcript related functions
// in the Rev.ai API.
type TranscriptService service

// Transcript represents a Rev.ai job transcript
type Transcript struct {
	Monologues []Monologue `json:"monologues"`
}

// Monologue represents a Rev.ai monologue
type Monologue struct {
	Speaker  int       `json:"speaker"`
	Elements []Element `json:"elements"`
}

// Element represents a Rev.ai element
type Element struct {
	Type       string  `json:"type"`
	Value      string  `json:"value"`
	Ts         float64 `json:"ts"`
	EndTs      float64 `json:"end_ts"`
	Confidence float64 `json:"confidence"`
}

// GetTranscriptParams specifies the parameters to the
// TranscriptService.Get method.
type GetTranscriptParams struct {
	JobID  string
	Accept string
}

// Get the developer's account information
// https://www.rev.ai/docs#operation/GetAccount
func (s *TranscriptService) Get(ctx context.Context, params *GetTranscriptParams) (*Transcript, error) {
	accept := params.Accept
	if accept != TextPlainAcceptHeader {
		accept = RevTranscriptJSONAcceptHeader
	}

	urlPath := "/speechtotext/v1/jobs/" + params.JobID + "/transcript"

	req, err := s.client.newRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	req.Header.Add("Accept", accept)

	var account Account
	if err := s.client.doJSON(ctx, req, &account); err != nil {
		return nil, err
	}

	return &account, nil
}
