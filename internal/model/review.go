package models

import "time"

type ReviewStatus string

const (
    StatusPending ReviewStatus = "pending"
    StatusReady   ReviewStatus = "ready"
    StatusFailed  ReviewStatus = "failed"
)

type Review struct {
    ID         string       `json:"id"`
    BookID     int          `json:"book_id"`
    Score      int          `json:"score"`
    ReviewText string       `json:"review_text"`
    Status     ReviewStatus `json:"status"`
    BookTitle  string       `json:"book_title,omitempty"`
    Authors    string       `json:"authors,omitempty"`
    CoverURL   string       `json:"cover_url,omitempty"`
    CreatedAt  time.Time    `json:"created_at"`
    UpdatedAt  time.Time    `json:"updated_at"`
}
