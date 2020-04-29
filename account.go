package revai

import (
	"context"
	"fmt"
	"net/http"
)

type AccountService service

type Account struct {
	Email          string `json:"email"`
	BalanceSeconds int    `json:"balance_seconds"`
}

func (s *AccountService) Get(ctx context.Context) (*Account, error) {
	urlPath := "/speechtotext/v1/account"

	req, err := s.client.newRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var account Account
	if err := s.client.do(ctx, req, &account); err != nil {
		return nil, err
	}

	return &account, nil
}
