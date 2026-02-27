package pocketsmith

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// PageInfo contains pagination metadata parsed from response headers.
type PageInfo struct {
	// Total is the total number of records across all pages.
	Total int
	// PerPage is the number of results per page.
	PerPage int
	// Page is the current page number, parsed from the request URL.
	Page int
	// Links contains pagination URLs keyed by relation: "first", "last", "next", "prev".
	Links map[string]string
}

// HasNext reports whether there is a next page of results.
func (p PageInfo) HasNext() bool {
	_, ok := p.Links["next"]
	return ok
}

// HasPrevious reports whether there is a previous page of results.
func (p PageInfo) HasPrevious() bool {
	_, ok := p.Links["prev"]
	return ok
}

// Last returns the total number of pages, derived from Total and PerPage.
// Returns 0 if PerPage is not set.
func (p PageInfo) Last() int {
	if p.PerPage <= 0 {
		return 0
	}
	return (p.Total + p.PerPage - 1) / p.PerPage
}

func (p PageInfo) String() string {
	return fmt.Sprintf("page %d (per_page=%d, total=%d)", p.Page, p.PerPage, p.Total)
}

// Page represents a single page of results from a paginated API endpoint.
type Page[T any] struct {
	Items    []T
	PageInfo PageInfo
}

// doList performs a GET request to the given URL and decodes the JSON response
// into a Page[T], including pagination metadata from response headers.
func doList[T any](ctx context.Context, c *Client, req *http.Request) (Page[T], error) {
	var page Page[T]

	resp, err := c.do(ctx, req)
	if err != nil {
		return page, err
	}

	page.PageInfo = parsePaginationHeaders(resp)

	if err := c.decodeJSON(resp, &page.Items); err != nil {
		return page, fmt.Errorf("decoding response: %w", err)
	}

	if page.Items == nil {
		page.Items = make([]T, 0)
	}

	return page, nil
}

// parsePaginationHeaders extracts PageInfo from the PocketSmith response headers:
// Total, Per-Page, and the RFC 5988 Link header.
func parsePaginationHeaders(resp *http.Response) PageInfo {
	info := PageInfo{
		Links: make(map[string]string),
	}

	if v := resp.Header.Get("Total"); v != "" {
		info.Total, _ = strconv.Atoi(v)
	}

	if v := resp.Header.Get("Per-Page"); v != "" {
		info.PerPage, _ = strconv.Atoi(v)
	}

	if v := resp.Header.Get("Link"); v != "" {
		info.Links = parseLinkHeader(v)
	}

	// Try to extract current page from the request URL.
	if resp.Request != nil {
		if p := resp.Request.URL.Query().Get("page"); p != "" {
			info.Page, _ = strconv.Atoi(p)
		} else {
			info.Page = 1
		}
	}

	return info
}

// parseLinkHeader parses an RFC 5988 Link header value into a map of
// relation type to URL. For example:
//
//	<https://api.example.com/items?page=2>; rel="next", <https://api.example.com/items?page=5>; rel="last"
//
// returns {"next": "https://api.example.com/items?page=2", "last": "https://api.example.com/items?page=5"}
func parseLinkHeader(header string) map[string]string {
	links := make(map[string]string)

	if header == "" {
		return links
	}

	// Split on commas to get individual link entries.
	// A link entry looks like: <URL>; rel="relation"
	for _, entry := range splitLinks(header) {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		// Split entry into parts on semicolons: <URL> ; rel="next" ; other="value"
		parts := strings.Split(entry, ";")
		if len(parts) < 2 {
			continue
		}

		// Extract URL from angle brackets.
		rawURL := strings.TrimSpace(parts[0])
		if !strings.HasPrefix(rawURL, "<") || !strings.HasSuffix(rawURL, ">") {
			continue
		}
		linkURL := rawURL[1 : len(rawURL)-1]

		// Find the rel parameter.
		for _, param := range parts[1:] {
			param = strings.TrimSpace(param)
			if !strings.HasPrefix(param, "rel=") {
				continue
			}
			rel := strings.TrimPrefix(param, "rel=")
			rel = strings.Trim(rel, `"`)
			if rel != "" {
				links[rel] = linkURL
			}
		}
	}

	return links
}

// splitLinks splits a Link header value on commas, respecting angle-bracket
// delimited URLs that may contain commas (though rare in practice).
func splitLinks(header string) []string {
	var entries []string
	var current strings.Builder
	depth := 0

	for _, r := range header {
		switch r {
		case '<':
			depth++
			current.WriteRune(r)
		case '>':
			depth--
			current.WriteRune(r)
		case ',':
			if depth == 0 {
				entries = append(entries, current.String())
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		entries = append(entries, current.String())
	}

	return entries
}

// PaginationParams contains the common pagination parameters supported by
// all paginated PocketSmith API endpoints.
type PaginationParams struct {
	// Page is the page number to retrieve. Pages are 1-indexed.
	// If zero, the API defaults to page 1.
	Page int
	// PerPage is the number of results per page.
	// The API default is 30, with a minimum of 10 and maximum of 1000.
	PerPage int
}

// EncodePage adds the pagination parameters to the given url.Values.
func (p PaginationParams) EncodePage(q url.Values) {
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	if p.PerPage > 0 {
		q.Set("per_page", strconv.Itoa(p.PerPage))
	}
}
