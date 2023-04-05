package revai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// TopicExtractionService provides access to topic extraction related functions
// in the Rev.ai API.
type TopicExtractionService service

// TopicExtraction represents a rev.ai asycn Topic Extraction.
type TopicExtraction struct {
	ID            string    `json:"id,omitempty"`
	Status        string    `json:"status,omitempty"`
	Language      string    `json:"language,omitempty"`
	CreatedOn     time.Time `json:"created_on,omitempty"`
	Type          string    `json:"type,omitempty"`
	Duration      string    `json:"processed_duration_seconds,omitempty"`
	WordCount     int       `json:"word_count,omitempty"`
	CallbackURL   string    `json:"callback_url,omitempty"`
	Metadata      string    `json:"metadata,omitempty"`
	CompletedOn   time.Time `json:"completed_on,omitempty"`
	Failure       string    `json:"title,omitempty"`
	FailureReason string    `json:"detail,omitempty"`
	AllowedValues []string  `json:"allowed_values,omitempty"`
	CurrentStatus string    `json:"current_value,omitempty"`
}

type TopicExtractionResults struct {
	Topics  []TopicExtractionTopics `json:"topics,omitempty"`
	Message string
}

type TopicExtractionTopics struct {
	Name       string                      `json:"topic_name"`
	Score      float64                     `json:"score"`
	Informants []TopicExtractionInformants `json:"informants"`
}

type TopicExtractionInformants struct {
	Content string `json:"content"`
	// Json Submission
	Start float64 `json:"ts,omitempty"`
	End   float64 `json:"end_ts,omitempty"`
	// PlainText Submission
	Offset int `json:"offset,omitempty"`
	Length int `json:"length,omitempty"`
}

type TopicExtractionPlainTextParams struct {
	CallbackURL        string     `json:"callback_url,omitempty"`
	Metadata           string     `json:"metadata,omitempty"`
	NotificationConfig *UrlConfig `json:"notification_config,omitempty"`
	DeleteSeconds      int        `json:"delete_after_seconds,omitempty"`
	Language           string     `json:"language,omitempty"`
	Text               string     `json:"text,omitempty"`
}

// SubmitFile starts an asynchronous job to extract topics from a transcript.
// https://www.rev.ai/docs#operation/SubmitTopicExtractionJob
func (s *TopicExtractionService) SubmitPlainText(ctx context.Context, params *TopicExtractionPlainTextParams) (*TopicExtraction, error) {
	if params.Text == "" {
		return nil, errors.New("text is required")
	}

	req, err := s.client.newRequest(http.MethodPost, "/topic_extraction/v1/jobs", params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	req.Header.Add("Content-Type", TextPlainHeader)

	var j TopicExtraction
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

type TopicExtractionJsonParams struct {
	Metadata           string     `json:"metadata,omitempty"`
	NotificationConfig *UrlConfig `json:"notification_config,omitempty"`
	DeleteSeconds      int        `json:"delete_after_seconds,omitempty"`
	Language           string     `json:"language,omitempty"`
	Transcript         Transcript `json:"json,omitempty"`
}

// SubmitFile starts an asynchronous job to extract topics from a transcript.
// https://www.rev.ai/docs#operation/SubmitTopicExtractionJob
func (s *TopicExtractionService) SubmitTranscriptJson(ctx context.Context, params *TopicExtractionJsonParams) (*TopicExtraction, error) {
	req, err := s.client.newRequest(http.MethodPost, "/topic_extraction/v1/jobs", params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var j TopicExtraction
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// GetTopicExtractionParams specifies the parameters to the
// TopicExtractionService.Get method.
type GetTopicExtractionParams struct {
	ID        string
	Threshold *float64
}

// Get returns the topic extraction output results.
// https://www.rev.ai/docs#tag/GetTopicExtractionResultById
func (s *TopicExtractionService) Get(ctx context.Context, params *GetTopicExtractionParams) (*TopicExtractionResults, error) {
	urlPath := "/topic_extraction/v1/jobs/" + params.ID + "/result"

	var threshold interface{}
	if params.Threshold != nil {
		threshold = *params.Threshold
	}

	req, err := s.client.newRequest(http.MethodGet, urlPath, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	req.Header.Set("Accept", RevTopicJSONHeader)

	var r TopicExtractionResults
	if err = s.client.doJSON(ctx, req, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// GetJobById returns the topic extraction output. Includes Status, failure reasons.
// https://www.rev.ai/docs#tag/GetTopicExtractionJobById
func (s *TopicExtractionService) GetJobById(ctx context.Context, params *GetTopicExtractionParams) (*TopicExtraction, error) {
	urlPath := "/topic_extraction/v1/jobs/" + params.ID

	req, err := s.client.newRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	req.Header.Add("Accept", RevTopicJSONHeader)

	var j TopicExtraction
	if err = s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// Delete deletes a topic extraction job
// https://www.rev.ai/docs#operation/DeleteTopicExtractionJobById
func (s *TopicExtractionService) Delete(ctx context.Context, params *DeleteParams) (*TopicExtraction, error) {
	if params.ID == "" {
		return nil, errors.New("id is required")
	}

	urlPath := "/topic_extraction/v1/jobs/" + params.ID

	req, err := s.client.newRequest(http.MethodDelete, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var j TopicExtraction
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// List gets a list of topic extraction jobs submitted within the last 30 days
// in reverse chronological order up to the provided limit number of jobs per call.
// https://www.rev.ai/docs#operation/GetListOfTopicExtractionJobs
func (s *TopicExtractionService) List(ctx context.Context, params *ListParams) ([]*TopicExtraction, error) {
	urlPath := "/topic_extraction/v1/jobs"

	req, err := s.client.newRequest(http.MethodGet, urlPath, params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var jobs []*TopicExtraction
	if err := s.client.doJSON(ctx, req, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}
