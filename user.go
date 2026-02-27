package pocketsmith

import (
	"context"
	"net/http"
	"time"
)

// User defines a PocketSmith user.
type User struct {
	ID                      int    `json:"id"`
	Login                   string `json:"login"`
	Name                    string `json:"name"`
	Email                   string `json:"email"`
	AvatarURL               string `json:"avatar_url"`
	BetaUser                bool   `json:"beta_user"`
	TimeZone                string `json:"time_zone"`
	WeekStartDay            int    `json:"week_start_day"`
	IsReviewingTransactions bool   `json:"is_reviewing_transactions"`
	BaseCurrencyCode        string `json:"base_currency_code"`
	AlwaysShowBaseCurrency  bool   `json:"always_show_base_currency"`
	UsingMultipleCurrencies bool   `json:"using_multiple_currencies"`

	AvailableAccounts int `json:"available_accounts"`
	AvailableBudgets  int `json:"available_budgets"`

	ForecastLastUpdatedAt    time.Time  `json:"forecast_last_updated_at"`
	ForecastLastAccessedAt   time.Time  `json:"forecast_last_accessed_at"`
	ForecastStartDate        CustomTime `json:"forecast_start_date"`
	ForecastEndDate          CustomTime `json:"forecast_end_date"`
	ForecastDeferRecalculate bool       `json:"forecast_defer_recalculate"`
	ForecastNeedsRecalculate bool       `json:"forecast_needs_recalculate"`

	LastLoggedInAt time.Time `json:"last_logged_in_at"`
	LastActivityAt time.Time `json:"last_activity_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Client) GetMe(ctx context.Context) (user User, err error) {
	var data User
	u := c.makeURL("/me")
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return data, err
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return data, err
	}
	if err := c.decodeJSON(resp, &data); err != nil {
		return data, err
	}
	return data, nil
}
