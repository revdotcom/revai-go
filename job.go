package revai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"time"
)

type JobService service

type NewJob struct {
	ID        string    `json:"id"`
	CreatedOn time.Time `json:"created_on"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	Metadata  string    `json:"metadata"`
}

type NewJobParams struct {
	Media      io.Reader
	Filename   string
	JobOptions *JobOptions
}

type JobOptions struct {
	SkipDiarization      bool   `json:"skip_diarization,omitempty"`
	SkipPunctuation      bool   `json:"skip_punctuation,omitempty"`
	RemoveDisfluencies   bool   `json:"remove_disfluencies,omitempty"`
	FilterProfanity      bool   `json:"filter_profanity,omitempty"`
	SpeakerChannelsCount int    `json:"selenaninphe@gmail.com,omitempty"`
	Metadata             string `json:"metadata,omitempty"`
	CallbackURL          string `json:"callback_url,omitempty"`
}

func (s *JobService) SubmitFile(ctx context.Context, params *NewJobParams) (*NewJob, error) {
	if params.Filename == "" {
		return nil, errors.New("filename is required")
	}

	if params.Media == nil {
		return nil, errors.New("media is required")
	}

	body := &bytes.Buffer{}

	mw := multipart.NewWriter(body)

	if err := makeReaderPart(mw, "media", params.Filename, params.Media); err != nil {
		return nil, err
	}

	if params.JobOptions != nil {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(params.JobOptions); err != nil {
			return nil, err
		}

		if err := makeStringPart(mw, "options", buf.String()); err != nil {
			return nil, err
		}
	}

	if err := mw.Close(); err != nil {
		return nil, fmt.Errorf("failed closing multipart-form writer %w", err)
	}

	req, err := s.client.newMultiPartRequest(mw, "/speechtotext/v1/jobs", body)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var newJob NewJob
	if err := s.client.do(ctx, req, &newJob); err != nil {
		return nil, err
	}

	return &newJob, nil
}
