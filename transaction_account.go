package pocketsmith

import "time"

// Account represents the account associated with a transaction.
type TransactionAccount struct {
	ID                           int64       `json:"id"`
	AccountID                    int64       `json:"account_id"`
	Name                         string      `json:"name"`
	LatestFeedName               string      `json:"latest_feed_name"`
	Number                       string      `json:"number"`
	Type                         string      `json:"type"`
	Offline                      bool        `json:"offline"`
	IsNetWorth                   bool        `json:"is_net_worth"`
	IncludeInNetWorth            bool        `json:"include_in_net_worth"`
	CurrencyCode                 string      `json:"currency_code"`
	CurrentBalance               float64     `json:"current_balance"`
	CurrentBalanceInBaseCurrency float64     `json:"current_balance_in_base_currency"`
	CurrentBalanceExchangeRate   *float64    `json:"current_balance_exchange_rate,omitempty"`
	CurrentBalanceDate           CustomTime  `json:"current_balance_date"`
	CurrentBalanceSource         string      `json:"current_balance_source"`
	DataFeedsBalanceType         string      `json:"data_feeds_balance_type"`
	SafeBalance                  *float64    `json:"safe_balance,omitempty"`
	SafeBalanceInBaseCurrency    *float64    `json:"safe_balance_in_base_currency,omitempty"`
	HasSafeBalanceAdjustment     bool        `json:"has_safe_balance_adjustment"`
	StartingBalance              float64     `json:"starting_balance"`
	StartingBalanceDate          CustomTime  `json:"starting_balance_date"`
	Institution                  Institution `json:"institution"`
	DataFeedsAccountID           string      `json:"data_feeds_account_id"`
	DataFeedsConnectionID        string      `json:"data_feeds_connection_id"`
	CreatedAt                    time.Time   `json:"created_at"`
	UpdatedAt                    time.Time   `json:"updated_at"`
}
