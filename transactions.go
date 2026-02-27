package pocketsmith

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type Transaction struct {
	ID                   int64              `json:"id"`
	Payee                string             `json:"payee"`
	OriginalPayee        string             `json:"original_payee"`
	Date                 CustomTime         `json:"date"`
	UploadSource         string             `json:"upload_source"`
	Category             *Category          `json:"category,omitempty"`
	ClosingBalance       float64            `json:"closing_balance"`
	ChequeNumber         *string            `json:"cheque_number,omitempty"`
	Memo                 *string            `json:"memo,omitempty"`
	Amount               float64            `json:"amount"`
	AmountInBaseCurrency float64            `json:"amount_in_base_currency"`
	Type                 string             `json:"type"`
	IsTransfer           bool               `json:"is_transfer"`
	NeedsReview          bool               `json:"needs_review"`
	Status               string             `json:"status"`
	Note                 *string            `json:"note,omitempty"`
	Labels               []string           `json:"labels,omitempty"`
	TransactionAccount   TransactionAccount `json:"transaction_account"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

type ListTransactionsParams struct {
	PaginationParams
	Uncategorized bool
}

func (l ListTransactionsParams) Values() url.Values {
	q := url.Values{}
	l.EncodePage(q)
	if l.Uncategorized {
		q.Add("uncategorised", "1")
	}
	return q
}

func (c *Client) ListTransactionsInUser(ctx context.Context, userID int, p ListTransactionsParams) (Page[Transaction], error) {
	u := c.makeURL(path.Join("users", strconv.Itoa(userID), "transactions"))
	u.RawQuery = p.Values().Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return Page[Transaction]{}, err
	}
	return doList[Transaction](ctx, c, req)
}

type UpdateTransactionParams struct {
	Memo         *string  `json:"memo,omitempty"`
	ChequeNumber *string  `json:"cheque_number,omitempty"`
	Payee        *string  `json:"payee,omitempty"`
	Amount       *float64 `json:"amount,omitempty"`
	Date         *string  `json:"date,omitempty"`
	IsTransfer   *bool    `json:"is_transfer,omitempty"`
	CategoryID   *int64   `json:"category_id,omitempty"`
	Note         *string  `json:"note,omitempty"`
	NeedsReview  *bool    `json:"needs_review,omitempty"`
	Labels       *string  `json:"labels,omitempty"`
}

func (c *Client) UpdateTransaction(ctx context.Context, id int64, p UpdateTransactionParams) (Transaction, error) {
	tx := Transaction{}
	jsonData, err := json.Marshal(p)
	if err != nil {
		return tx, err
	}

	u := c.makeURL(path.Join("transactions", fmt.Sprintf("%d", id)))
	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return tx, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.do(ctx, req)
	if err != nil {
		return tx, err
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
		return tx, err
	}

	return tx, nil
}
