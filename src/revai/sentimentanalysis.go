package revai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	PositiveSentiment = "positive"
	NegativeSentiment = "negative"
	NeutralSentiment  = "neutral"
)

// SentimentAnalysisService provides access to languageid related functions
// in the Rev.ai API.
type SentimentAnalysisService service

// Job represents a rev.ai asycn job.
type SentimentAnalysis struct {
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

type SentimentAnalysisResults struct {
	Messages []SentimentAnalysisMessages `json:"messages,omitempty"`
}

type SentimentAnalysisMessages struct {
	Sentiment string  `json:"sentiment"`
	Content   string  `json:"content"`
	Score     float64 `json:"score"`
	// Json Submission
	Start float64 `json:"ts,omitempty"`
	End   float64 `json:"end_ts,omitempty"`
	// PlainText Submission
	Offset int `json:"offset,omitempty"`
	Length int `json:"length,omitempty"`
}

type SentimentAnalysisPlainTextParams struct {
	CallbackURL        *string    `json:"callback_url,omitempty"`
	Metadata           string     `json:"metadata,omitempty"`
	NotificationConfig *UrlConfig `json:"notification_config,omitempty"`
	DeleteSeconds      int        `json:"delete_after_seconds,omitempty"`
	Language           string     `json:"language,omitempty"`
	Text               string     `json:"text,omitempty"`
}

// SubmitFile starts an asynchronous job to extract topics from a transcript.
// https://www.rev.ai/docs#operation/SubmitSentimentAnalysisJob
func (s *SentimentAnalysisService) SubmitPlainText(ctx context.Context, params *SentimentAnalysisPlainTextParams) (*SentimentAnalysis, error) {
	if params.Text == "" {
		return nil, errors.New("text is required")
	}

	req, err := s.client.newRequest(http.MethodPost, "/sentiment_analysis/v1/jobs", params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	req.Header.Add("Content-Type", TextPlainHeader)

	var j SentimentAnalysis
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

type SentimentAnalysisJsonParams struct {
	Metadata           string     `json:"metadata,omitempty"`
	NotificationConfig *UrlConfig `json:"notification_config,omitempty"`
	DeleteSeconds      int        `json:"delete_after_seconds,omitempty"`
	Language           string     `json:"language,omitempty"`
	Transcript         Transcript `json:"json,omitempty"`
}

// SubmitFile starts an asynchronous job to extract topics from a transcript.
// https://www.rev.ai/docs#operation/SubmitSentimentAnalysisJob
func (s *SentimentAnalysisService) SubmitTranscriptJson(ctx context.Context, params *SentimentAnalysisJsonParams) (*SentimentAnalysis, error) {
	req, err := s.client.newRequest(http.MethodPost, "/sentiment_analysis/v1/jobs", params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var j SentimentAnalysis
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// GetSentimentAnalysisParams specifies the parameters to the
// SentimentAnalysisService.Get method.
type GetSentimentAnalysisParams struct {
	ID     string
	Filter *string
}

// Get returns the sentiment analysis output results.
// https://www.rev.ai/docs#tag/GetSentimentAnalysisResultById
func (s *SentimentAnalysisService) Get(ctx context.Context, params *GetSentimentAnalysisParams) (*SentimentAnalysisResults, error) {
	urlPath := "/sentiment_analysis/v1/jobs/" + params.ID + "/result"

	var filter interface{}
	if params.Filter != nil {
		filter = *params.Filter
	}

	req, err := s.client.newRequest(http.MethodGet, urlPath, filter)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	req.Header.Set("Accept", RevSentimentJSONHeader)

	var r SentimentAnalysisResults
	if err = s.client.doJSON(ctx, req, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// GetJobById returns the topic extraction output. Includes Status, failure reasons.
// https://www.rev.ai/docs#tag/GetSentimentAnalysisJobById
func (s *SentimentAnalysisService) GetJobById(ctx context.Context, params *GetSentimentAnalysisParams) (*SentimentAnalysis, error) {
	urlPath := "/sentiment_analysis/v1/jobs/" + params.ID

	req, err := s.client.newRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	req.Header.Set("Accept", RevSentimentJSONHeader)

	var j SentimentAnalysis
	if err = s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// Delete deletes a topic extraction job
// https://www.rev.ai/docs#operation/DeleteSentimentAnalysisJobById
func (s *SentimentAnalysisService) Delete(ctx context.Context, params *DeleteParams) (*SentimentAnalysis, error) {
	if params.ID == "" {
		return nil, errors.New("id is required")
	}

	urlPath := "/sentiment_analysis/v1/jobs/" + params.ID

	req, err := s.client.newRequest(http.MethodDelete, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var j SentimentAnalysis
	if err := s.client.doJSON(ctx, req, &j); err != nil {
		return nil, err
	}

	return &j, nil
}

// List gets a list of topic extraction jobs submitted within the last 30 days
// in reverse chronological order up to the provided limit number of jobs per call.
// https://www.rev.ai/docs#operation/GetListOfSentimentAnalysisJobs
func (s *SentimentAnalysisService) List(ctx context.Context, params *ListParams) ([]*SentimentAnalysis, error) {
	urlPath := "/sentiment_analysis/v1/jobs"

	req, err := s.client.newRequest(http.MethodGet, urlPath, params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var jobs []*SentimentAnalysis
	if err := s.client.doJSON(ctx, req, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}
