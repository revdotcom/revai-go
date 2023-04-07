package revai

import (
	"context"
	"fmt"
	"net/http"
)

// AccountService provides access to the account related functions
// in the Rev.ai API.
type AccountService service

// Account is the developer's account information
type Account struct {
	Email            string  `json:"email"`
	FreeBalance      float64 `json:"free_balance"`
	PurchasedBalance float64 `json:"purchased_balance"`
	TotalBalance     float64 `json:"total_balance"`
	InvoicedBalance  float64 `json:"invoiced_balance"`
	BalanceSeconds   int     `json:"balance_seconds"`
	HipaaEnabled     bool    `json:"hipaa_enabled"`
}

// Get the developer's account information
// https://www.rev.ai/docs#operation/GetAccount
func (s *AccountService) Get(ctx context.Context) (*Account, error) {
	urlPath := "/speechtotext/v1/account"

	req, err := s.client.newRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var account Account
	if err := s.client.doJSON(ctx, req, &account); err != nil {
		return nil, err
	}

	return &account, nil
}
