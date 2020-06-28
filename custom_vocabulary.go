package revai

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// CustomVocabularyService provides access to the custom vocabulary related functions
// in the Rev.ai API.
type CustomVocabularyService service

// CustomVocabulary represents a Rev.ai custom vocabulary
type CustomVocabulary struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	CreatedOn     time.Time `json:"created_on"`
	CompletedOn   time.Time `json:"completed_on"`
	CallbackURL   string    `json:"callback_url"`
	Failure       string    `json:"failure"`
	FailureDetail string    `json:"failure_detail"`
}

// CreateCustomVocabularyParams specifies the parameters to the
// CustomVocabularyService.Create method.
type CreateCustomVocabularyParams struct {
	CustomVocabularies []Phrase `json:"custom_vocabularies"`
	Metadata           string   `json:"metadata,omitempty"`
	CallbackURL        string   `json:"callback_url,omitempty"`
}

type Phrase struct {
	Phrases []string `json:"phrases"`
}

// Create submits a Custom Vocabulary for asynchronous processing.
// https://www.rev.ai/docs/streaming#operation/SubmitCustomVocabulary
func (s *CustomVocabularyService) Create(ctx context.Context, params *CreateCustomVocabularyParams) (*CustomVocabulary, error) {
	urlPath := "/speechtotext/v1/vocabularies"

	req, err := s.client.newRequest(http.MethodPost, urlPath, params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var vocabulary CustomVocabulary
	if err := s.client.doJSON(ctx, req, &vocabulary); err != nil {
		return nil, err
	}

	return &vocabulary, nil
}

// GetCustomVocabularyParams specifies the parameters to the
// CustomVocabularyService.Get method.
type GetCustomVocabularyParams struct {
	ID string
}

// Get gets the custom vocabulary processing information
// https://www.rev.ai/docs/streaming#operation/GetCustomVocabulary
func (s *CustomVocabularyService) Get(ctx context.Context, params *GetCustomVocabularyParams) (*CustomVocabulary, error) {
	urlPath := "/speechtotext/v1/vocabularies/" + params.ID

	req, err := s.client.newRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var vocabulary CustomVocabulary
	if err := s.client.doJSON(ctx, req, &vocabulary); err != nil {
		return nil, err
	}

	return &vocabulary, nil
}

// ListCustomVocabularyParams specifies the parameters to the
// CustomVocabularyService.List method.
type ListCustomVocabularyParams struct {
	Limit int `url:"limit,omitempty"`
}

// List gets a list of most recent custom vocabularies' processing information
// https://www.rev.ai/docs/streaming#operation/GetCustomVocabularies
func (s *CustomVocabularyService) List(ctx context.Context, params *ListCustomVocabularyParams) ([]*CustomVocabulary, error) {
	urlPath := "/speechtotext/v1/vocabularies"

	req, err := s.client.newRequest(http.MethodGet, urlPath, params)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var vocabularies []*CustomVocabulary
	if err := s.client.doJSON(ctx, req, &vocabularies); err != nil {
		return nil, err
	}

	return vocabularies, nil
}

// DeleteCustomVocabularyParams specifies the parameters to the
// CustomVocabularyService.Delete method.
type DeleteCustomVocabularyParams struct {
	ID string
}

// Delete deletes the custom vocabulary.
// https://www.rev.ai/docs/streaming#operation/DeleteCustomVocabulary
func (s *CustomVocabularyService) Delete(ctx context.Context, params *DeleteCustomVocabularyParams) error {
	urlPath := "/speechtotext/v1/vocabularies/" + params.ID

	req, err := s.client.newRequest(http.MethodDelete, urlPath, nil)
	if err != nil {
		return fmt.Errorf("failed creating request %w", err)
	}

	if err := s.client.doJSON(ctx, req, nil); err != nil {
		return err
	}

	return nil
}
