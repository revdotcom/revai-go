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

// JobService provides access to the jobs related functions
// in the Rev.ai API.
type JobService service

// Job represents a rev.ai asycn job.
type Job struct {
	ID              string       `json:"id"`
	CreatedOn       time.Time    `json:"created_on"`
	Name            string       `json:"name"`
	Status          string       `json:"status"`
	Type            string       `json:"type"`
	Language        string       `json:"language"`
	Metadata        string       `json:"metadata,omitempty"`
	CompletedOn     time.Time    `json:"completed_on,omitempty"`
	CallbackURL     string       `json:"callback_url,omitempty"`
	DurationSeconds float32      `json:"duration_seconds,omitempty"`
	MediaURL        string       `json:"media_url,omitempty"`
	FailureParam    JobFailParam `json:"parameter,omitempty"`
	FailureDetail   string       `json:"error,omitempty"`
}

type JobFailParam struct {
	FailureDetail []string `json:"media_url,omitempty"`
}

// NewFileJobParams specifies the parameters to the
// JobService.SubmitFile method.
type NewFileJobParams struct {
	Media      io.Reader
	Filename   string
	JobOptions *JobOptions
}

// JobOptions specifies the options to the
// JobService.SubmitFile method.
type JobOptions struct {
	Language             string                      `json:"language,omitempty"`
	NotificationConfig   *UrlConfig                  `json:"notification_config,omitempty"`
	CallbackURL          string                      `json:"callback_url,omitempty"`
	SkipDiarization      bool                        `json:"skip_diarization,omitempty"`
	SkipPunctuation      bool                        `json:"skip_punctuation,omitempty"`
	SkipPostProcessing   bool                        `json:"skip_postprocessing,omitempty"`
	RemoveDisfluencies   bool                        `json:"remove_disfluencies,omitempty"`
	RemoveAtmospherics   bool                        `json:"remove_atmospherics,omitempty"`
	FilterProfanity      bool                        `json:"filter_profanity,omitempty"`
	SpeakerChannelsCount int                         `json:"speaker_channels_count,omitempty"`
	Metadata             string                      `json:"metadata,omitempty"`
	CustomVocabularyId   string                      `json:"custom_vocabulary_id,omitempty"`
	CustomVocabularies   []JobOptionCustomVocabulary `json:"custom_vocabularies,omitempty"`
	Transcriber          string                      `json:"transcriber,omitempty"`
	Verbatim             bool                        `json:"verbatim,omitempty"`
}

type JobOptionCustomVocabulary struct {
	Phrases []string `json:"phrases"`
}

// SubmitFile starts an asynchronous job to transcribe speech-to-text for a media file.
// https://www.rev.ai/docs#operation/SubmitTranscriptionJob
func (s *JobService) SubmitFile(ctx context.Context, params *NewFileJobParams) (*Job, error) {
	if params.Filename == "" {
		return nil, errors.New("filename is required")
	}

	if params.Media == nil {
		return nil, errors.New("media is required")
	}

	pr, pw := io.Pipe()

	mw := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		if err := makeReaderPart(mw, "media", params.Filename, params.Media); err != nil {
			pw.CloseWithError(err)
			return
		}

		if params.JobOptions != nil {
			buf := new(bytes.Buffer)
			if err := json.NewEncoder(buf).Encode(params.JobOptions); err != nil {
				pw.CloseWithError(err)
				return
			}

			if err := makeStringPart(mw, "options", buf.String()); err != nil {
				pw.CloseWithError(err)
				return
			}
		}

		if err := mw.Close(); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	req, err := s.client.newMultiPartRequest(mw, "/speechtotext/v1/jobs", pr)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var j Job
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// NewURLJobParams specifies the parameters to the
// JobService.SubmitURL method.
type NewURLJobParams struct {
	Language             string                      `json:"language,omitempty"`
	MediaURL             *string                     `json:"media_url,omitempty"`
	SourceConfig         *UrlConfig                  `json:"source_config,omitempty"`
	NotificationConfig   *UrlConfig                  `json:"notification_config,omitempty"`
	SkipDiarization      bool                        `json:"skip_diarization,omitempty"`
	SkipPunctuation      bool                        `json:"skip_punctuation,omitempty"`
	RemoveDisfluencies   bool                        `json:"remove_disfluencies,omitempty"`
	RemoveAtmospherics   bool                        `json:"remove_atmospherics,omitempty"`
	FilterProfanity      bool                        `json:"filter_profanity,omitempty"`
	SpeakerChannelsCount int                         `json:"speaker_channels_count,omitempty"`
	Metadata             string                      `json:"metadata,omitempty"`
	CallbackURL          string                      `json:"callback_url,omitempty"`
	CustomVocabularyId   string                      `json:"custom_vocabulary_id,omitempty"`
	CustomVocabularies   []JobOptionCustomVocabulary `json:"custom_vocabularies,omitempty"`
	DeleteSeconds        int                         `json:"delete_after_seconds,omitempty"`
	Transcriber          string                      `json:"transcriber,omitempty"`
	Verbatim             bool                        `json:"verbatim,omitempty"`
}

type UrlConfig struct {
	Url         string            `json:"url,omitempty"`
	AuthHeaders map[string]string `json:"auth_headers,omitempty"`
}

// SubmitURL starts an asynchronous job to transcribe speech-to-text for a media file.
// https://www.rev.ai/docs#operation/SubmitTranscriptionJob
func (s *JobService) SubmitURL(ctx context.Context, params *NewURLJobParams) (*Job, error) {
	if params.SourceConfig.Url == "" {
		return nil, errors.New("url is required")
	}

	req, err := s.client.newRequest(http.MethodPost, "/speechtotext/v1/jobs", params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var j Job
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// GetJobParams specifies the parameters to the
// JobService.Get method.
type GetJobParams struct {
	ID string
}

// Get returns information about a transcription job
// https://www.rev.ai/docs#operation/GetJobById
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
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// DeleteJobParams specifies the parameters to the
// JobService.Delete method.
type DeleteParams struct {
	ID string
}

// Delete deletes a transcription job
// https://www.rev.ai/docs#operation/DeleteJobById
func (s *JobService) Delete(ctx context.Context, params *DeleteParams) (*Job, error) {
	if params.ID == "" {
		return nil, errors.New("job id is required")
	}

	urlPath := "/speechtotext/v1/jobs/" + params.ID

	req, err := s.client.newRequest(http.MethodDelete, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var j Job
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return nil, nil
}

// ListParams specifies the optional query parameters to the
// Most List methods.
type ListParams struct {
	Limit         int    `url:"limit,omitempty"`
	StartingAfter string `url:"starting_after,omitempty"`
}

// List gets a list of transcription jobs submitted within the last 30 days
// in reverse chronological order up to the provided limit number of jobs per call.
// https://www.rev.ai/docs#operation/GetListOfJobs
func (s *JobService) List(ctx context.Context, params *ListParams) ([]*Job, error) {
	urlPath := "/speechtotext/v1/jobs"

	req, err := s.client.newRequest(http.MethodGet, urlPath, params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var jobs []*Job
	if err := s.client.doJSON(ctx, req, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}
