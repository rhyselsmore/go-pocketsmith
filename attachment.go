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

// Attachment represents a file attachment in PocketSmith, such as a receipt or screenshot.
type Attachment struct {
	ID              int64                     `json:"id"`
	Title           string                    `json:"title"`
	FileName        string                    `json:"file_name"`
	Type            string                    `json:"type"`
	ContentType     string                    `json:"content_type"`
	ContentTypeMeta AttachmentContentTypeMeta `json:"content_type_meta"`
	OriginalURL     string                    `json:"original_url"`
	Variants        AttachmentVariants        `json:"variants"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

// AttachmentContentTypeMeta contains metadata about an attachment's content type.
type AttachmentContentTypeMeta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Extension   string `json:"extension"`
}

// AttachmentVariants contains URLs for different sized versions of an attachment.
type AttachmentVariants struct {
	ThumbURL string `json:"thumb_url"`
	LargeURL string `json:"large_url"`
}

// ListAttachmentsInUserParams contains the parameters for listing a user's attachments.
type ListAttachmentsInUserParams struct {
	Unassigned bool
}

// Values returns the query parameters for the request.
func (l ListAttachmentsInUserParams) Values() url.Values {
	q := url.Values{}
	if l.Unassigned {
		q.Add("unassigned", "1")
	}
	return q
}

// ListAttachmentsInUser returns all attachments belonging to the given user.
// If Unassigned is set, only attachments not assigned to a transaction are returned.
func (c *Client) ListAttachmentsInUser(ctx context.Context, userID int, p ListAttachmentsInUserParams) ([]Attachment, error) {
	var att []Attachment
	u := c.makeURL(path.Join("users", strconv.Itoa(userID), "attachments"))
	u.RawQuery = p.Values().Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return att, err
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return att, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&att); err != nil {
		return att, fmt.Errorf("pocketsmith: json decode error: %w", err)
	}
	return att, nil
}

// CreateAttachmentInUserParams contains the parameters for creating an attachment.
// FileData should be a base64-encoded string of the file contents.
type CreateAttachmentInUserParams struct {
	Title    *string `json:"title,omitempty"`
	FileName string  `json:"file_name"`
	FileData string  `json:"file_data"`
}

// CreateAttachmentInUser uploads a new attachment for the given user.
func (c *Client) CreateAttachmentInUser(ctx context.Context, userId int64, p CreateAttachmentInUserParams) (Attachment, error) {
	att := Attachment{}
	jsonData, err := json.Marshal(p)
	if err != nil {
		return att, err
	}

	u := c.makeURL(path.Join("users", fmt.Sprintf("%d", userId), "attachments"))
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return att, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.do(ctx, req)
	if err != nil {
		return att, err
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&att); err != nil {
		return att, err
	}

	return att, nil
}
