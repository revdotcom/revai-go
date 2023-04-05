package revai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// LanguageIdService provides access to languageid related functions
// in the Rev.ai API.
type LanguageIdService service

// LanguageId represents a rev.ai asycn LanguageId.
type LanguageId struct {
	ID            string                `json:"id,omitempty"`
	CreatedOn     time.Time             `json:"created_on,omitempty"`
	Status        string                `json:"status,omitempty"`
	Type          string                `json:"type,omitempty"`
	Duration      float64               `json:"processed_duration_seconds,omitempty"`
	Metadata      string                `json:"metadata,omitempty"`
	CompletedOn   time.Time             `json:"completed_on,omitempty"`
	MediaURL      string                `json:"media_url,omitempty"`
	Failure       string                `json:"title,omitempty"`
	FailureReason string                `json:"detail,omitempty"`
	CurrentStatus string                `json:"current_value,omitempty"`
	AllowedValues []string              `json:"allowed_values,omitempty"`
	TopLanguage   string                `json:"top_language,omitempty"`
	Confidences   []LanguageConfidences `json:"language_confidences,omitempty"`
}

type LanguageConfidences struct {
	Language   string  `json:"language"`
	Confidence float64 `json:"confidence"`
}

// LanguageIdFileParams specifies the parameters to the
// LanguageIdService.SubmitFile method.
type LanguageIdFileParams struct {
	Options  *LanguageIdFileOptions
	Media    io.Reader
	Filename string
}

type LanguageIdFileOptions struct {
	NotificationConfig *UrlConfig `json:"notification_config,omitempty"`
	Metadata           string     `json:"metadata,omitempty"`
	DeleteSeconds      int        `json:"delete_after_seconds,omitempty"`
}

// SubmitFile starts an asynchronous job to get the language code for a media file.
// https://www.rev.ai/docs#operation/SubmitLanguageIdentificationJob
func (s *LanguageIdService) SubmitFile(ctx context.Context, params *LanguageIdFileParams) (*LanguageId, error) {
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

		if err := mw.Close(); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	req, err := s.client.newMultiPartRequest(mw, "/languageid/v1/jobs", pr)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var j LanguageId
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// LanguageIdParams specifies the parameters to the
// LanguageIdService.SubmitURL method.
type LanguageIdUrlParams struct {
	SourceConfig       *UrlConfig `json:"source_config,omitempty"`
	NotificationConfig *UrlConfig `json:"notification_config,omitempty"`
	Metadata           string     `json:"metadata,omitempty"`
	DeleteSeconds      int        `json:"delete_after_seconds,omitempty"`
}

// SubmitURL starts an asynchronous job to transcribe speech-to-text for a media file.
// https://www.rev.ai/docs#operation/SubmitLanguageIdentificationJob
func (s *LanguageIdService) SubmitURL(ctx context.Context, params *LanguageIdUrlParams) (*LanguageId, error) {
	if params.SourceConfig.Url == "" {
		return nil, errors.New("url is required")
	}

	req, err := s.client.newRequest(http.MethodPost, "/languageid/v1/jobs", params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var j LanguageId
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// GetLanguageIdParams specifies the parameters to the
// LanguageIdService.Get method.
type GetLanguageIdParams struct {
	ID string
}

// Get returns the languageid output results.
// https://www.rev.ai/docs#tag/GetLanguageIdentificationResultById
func (s *LanguageIdService) Get(ctx context.Context, params *GetLanguageIdParams) (*LanguageId, error) {
	urlPath := "/languageid/v1/jobs/" + params.ID + "/result"

	req, err := s.client.newRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	req.Header.Set("Accept", RevLanguageIdJSONHeader)

	var j LanguageId
	if err = s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// GetJobById returns the languageid output. Includes Status, failure reasons.
// https://www.rev.ai/docs#tag/GetLanguageIdentificationJobById
func (s *LanguageIdService) GetJobById(ctx context.Context, params *GetLanguageIdParams) (*LanguageId, error) {
	urlPath := "/languageid/v1/jobs/" + params.ID

	req, err := s.client.newRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	req.Header.Set("Accept", RevLanguageIdJSONHeader)

	var j LanguageId
	if err = s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// Delete deletes a language id job
// https://www.rev.ai/docs#operation/DeleteLanguageIdentificationJobById
func (s *LanguageIdService) Delete(ctx context.Context, params *DeleteParams) (*LanguageId, error) {
	if params.ID == "" {
		return nil, errors.New("job id is required")
	}

	urlPath := "/languageid/v1/jobs/" + params.ID

	req, err := s.client.newRequest(http.MethodDelete, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var j LanguageId
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// List gets a list of language id jobs submitted within the last 30 days
// in reverse chronological order up to the provided limit number of jobs per call.
// https://www.rev.ai/docs#operation/GetListOfLanguageIdentificationJobs
func (s *LanguageIdService) List(ctx context.Context, params *ListParams) ([]*LanguageId, error) {
	urlPath := "/languageid/v1/jobs"

	req, err := s.client.newRequest(http.MethodGet, urlPath, params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var jobs []*LanguageId
	if err := s.client.doJSON(ctx, req, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}
