package revai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type JobService service

type Job struct {
	ID              string    `json:"id"`
	CreatedOn       time.Time `json:"created_on"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	Type            string    `json:"type"`
	Metadata        string    `json:"metadata,omitempty"`
	CompletedOn     time.Time `json:"completed_on,omitempty"`
	CallbackURL     string    `json:"callback_url,omitempty"`
	DurationSeconds float32   `json:"duration_seconds,omitempty"`
	MediaURL        string    `json:"media_url,omitempty"`
	Failure         string    `json:"failure,omitempty"`
	FailureDetail   string    `json:"failure_detail,omitempty"`
}

type NewFileJobParams struct {
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

func (s *JobService) SubmitFile(ctx context.Context, params *NewFileJobParams) (*Job, error) {
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

	var j Job
	if err := s.client.do(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

type NewJobParams struct {
	MediaURL             string `json:"media_url"`
	SkipDiarization      bool   `json:"skip_diarization,omitempty"`
	SkipPunctuation      bool   `json:"skip_punctuation,omitempty"`
	RemoveDisfluencies   bool   `json:"remove_disfluencies,omitempty"`
	FilterProfanity      bool   `json:"filter_profanity,omitempty"`
	SpeakerChannelsCount int    `json:"selenaninphe@gmail.com,omitempty"`
	Metadata             string `json:"metadata,omitempty"`
	CallbackURL          string `json:"callback_url,omitempty"`
}

func (s *JobService) Submit(ctx context.Context, params *NewJobParams) (*Job, error) {
	if params.MediaURL == "" {
		return nil, errors.New("media url is required")
	}

	req, err := s.client.newRequest(http.MethodPost, "/speechtotext/v1/jobs", params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var j Job
	if err := s.client.do(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

type GetJobParams struct {
	ID string
}

func (s *JobService) Get(ctx context.Context, params *GetJobParams) (*Job, error) {
	if params.ID == "" {
		return nil, errors.New("job id is required")
	}

	urlPath := "/speechtotext/v1/jobs/" + params.ID

	req, err := s.client.newRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var j Job
	if err := s.client.do(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}
