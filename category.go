package pocketsmith

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type Category struct {
	ID              int64      `json:"id"`
	Title           string     `json:"title"`
	Colour          string     `json:"colour"`
	IsTransfer      bool       `json:"is_transfer"`
	IsBill          bool       `json:"is_bill"`
	RefundBehaviour *string    `json:"refund_behaviour"`
	Children        []Category `json:"children"`
	ParentID        *int64     `json:"parent_id"`
	RollUp          bool       `json:"roll_up"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type ListCategoriesParams struct {
	PaginationParams
	Uncategorized bool
}

func (l ListCategoriesParams) Values() url.Values {
	q := url.Values{}
	l.EncodePage(q)
	if l.Uncategorized {
		q.Add("uncategorised", "1")
	}
	return q
}

func (c *Client) ListCategories(ctx context.Context, userID int, p ListCategoriesParams) (Page[Category], error) {
	u := c.makeURL(path.Join("users", strconv.Itoa(userID), "categories"))
	u.RawQuery = p.Values().Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return Page[Category]{}, err
	}
	return doList[Category](ctx, c, req)
}
