package revai

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type CustomVocabularyService service

type CustomVocabulary struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	CreatedOn     time.Time `json:"created_on"`
	CompletedOn   time.Time `json:"completed_on"`
	CallbackURL   string    `json:"callback_url"`
	Failure       string    `json:"failure"`
	FailureDetail string    `json:"failure_detail"`
}

type GetCustomVocabularyParams struct {
	ID string
}

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

type ListCustomVocabularyParams struct {
	Limit int `url:"limit,omitempty"`
}

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
