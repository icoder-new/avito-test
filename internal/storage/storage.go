package storage

import (
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/models"
	"github.com/google/uuid"
)

type IStorage interface {
	CloseDB()

	Tender() ITender
	TenderHistory() ITenderHistory
	Organization() IOrganization
	Employee() IEmployee
	Bid() IBid
	BidHistory() IBidHistory
	BidFeedbacks() IBidFeedbacks
}

type ITender interface {
	GetAllTenders(limit, offset int, serviceType string) ([]models.Tender, error)
	CreateTender(tender models.Tender) (models.Tender, error)
	GetMyTenders(limit, offset int, userId, organizationId uuid.UUID) ([]models.Tender, error)
	GetTenderStatusById(tenderId uuid.UUID, username string) (string, error)
	CheckID(tenderId uuid.UUID) (bool, error)
	GetTenderById(tenderId uuid.UUID) (models.Tender, error)
	UpdateTender(tender models.Tender) error
}

type IOrganization interface {
	CheckID(id uuid.UUID) (bool, error)
	GetIdByUserId(userId uuid.UUID) (uuid.UUID, error)
}

type IEmployee interface {
	GetIdByUsername(username string) (uuid.UUID, error)
	IsOwnerOrganization(id, organizationId uuid.UUID) (bool, error)
	IsTenderOwner(id, tenderId uuid.UUID) (bool, error)
}

type ITenderHistory interface {
	CreateTenderHistory(id uuid.UUID, tender models.Tender) error
	GetVersionById(tenderId uuid.UUID) (int, error)
	GetTender(version int, userId, organizationId uuid.UUID) (models.Tender, error)
}

type IBid interface {
	GetBidsByFiltering(limit, offset int, authorId uuid.UUID) ([]models.Bid, error)
	GetBidsByUsername(authorId uuid.UUID) ([]models.Bid, error)
	CreateBid(bid models.Bid) (*models.Bid, error)
	GetBidListByTenderId(tenderId string, authorId uuid.UUID) ([]models.Bid, error)
	GetBidStatus(bidId string, authorId uuid.UUID) (status string, err error)
	ChangeBidStatus(bid *models.Bid) (*models.Bid, error)
	ChangeBidById(bid *models.Bid) (*models.Bid, error)
	SendDesicionByBid(bid *models.Bid) (*models.Bid, error)
	GetVersionBidById(bidId string) (version int, err error)
	RollbackBidById(bidId, username string, version int) (*models.Bid, error)
}

type IBidHistory interface {
	SaveBidHistory(bid models.Bid) error
}

type IBidFeedbacks interface {
	SendFeedbackByBid(bid *models.Bid, bidFeedback string) (*models.Bid, error)
	GetFeedbacksByTenderAndAuthor(tenderID string, authorId, requesterId uuid.UUID, limit int, offset int) ([]*models.BidFeedback, error)
}
