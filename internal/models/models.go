package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Tender struct {
	ID             uuid.UUID `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	ServiceType    string    `db:"service_type"`
	Status         string    `db:"status"`
	OrganizationID uuid.UUID `db:"organization_id"`
	ResponsibleID  uuid.UUID `db:"responsible_id"`
	Version        int       `db:"version"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type Bid struct {
	ID          uuid.UUID      `db:"id"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	Status      string         `db:"status"`
	TenderID    uuid.UUID      `db:"tender_id"`
	AuthorType  string         `db:"author_type"`
	AuthorID    uuid.UUID      `db:"author_id"`
	Version     int            `db:"version"`
	Decision    sql.NullString `db:"decision"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

type BidFeedback struct {
	ID        uuid.UUID `json:"id"`
	BidID     uuid.UUID `json:"bid_id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Feedback  string    `json:"feedback"`
	CreatedAt time.Time `json:"created_at"`
}
