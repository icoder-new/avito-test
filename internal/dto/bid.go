package dto

import (
	"time"

	"github.com/google/uuid"
)

type Decision string

const (
	Approved Decision = "Approved"
	Rejected Decision = "Rejected"
)

type BidAuthorType string

const (
	AuthorTypeOrganization BidAuthorType = "Organization"
	AuthorTypeUser         BidAuthorType = "User"
)

type BidRequest struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	TenderId    uuid.UUID     `json:"tenderId"`
	AuthorType  BidAuthorType `json:"authorType" validate:"required,oneof=Organization User"`
	AuthorId    uuid.UUID     `json:"authorId"`
}

type BidResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	AuthorType string    `json:"authorType"`
	AuthorID   uuid.UUID `json:"authorID"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"createdAt"`
}

type BidStatusRequest struct {
	BidID    string `validate:"required,uuid"`
	Status   Status `validate:"required,oneof=Created Published Canceled"`
	Username string `validate:"required"`
}

type ChangeBidRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	BidID       string `validate:"required,uuid"`
	Username    string `validate:"required"`
}

type SendBidDesicionRequest struct {
	BidID    string   `validate:"required,uuid"`
	Decision Decision `validate:"required,oneof=Approved Rejected"`
	Username string   `validate:"required"`
}

type SendBidFeedbackRequest struct {
	BidID       string `validate:"required,uuid"`
	BidFeedback string `validate:"required"`
	Username    string `validate:"required"`
}

type BidFeedbackResponse struct {
	ID        uuid.UUID `json:"id"`
	BidID     uuid.UUID `json:"bid_id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Feedback  string    `json:"feedback"`
	CreatedAt time.Time `json:"created_at"`
}
