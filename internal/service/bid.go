package service

import (
	"database/sql"
	"fmt"
	"time"

	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/config"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/dto"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/models"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/storage"
	customerrors "git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/errors"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/logger"
	"github.com/google/uuid"
)

type IBid interface {
	GetBidsByFiltering(limit, offset int, username string) ([]dto.BidResponse, error)
	GetBidsByUsername(username string) ([]dto.BidResponse, error)
	CreateBid(bidRequest dto.BidRequest) (*dto.BidResponse, error)
	GetBidListByTenderId(tenderId, username string) ([]dto.BidResponse, error)
	GetBidStatus(bidId, username string) (status string, err error)
	ChangeBidStatus(bidReq *dto.BidStatusRequest) (*dto.BidResponse, error)
	ChangeBidById(bidReq *dto.ChangeBidRequest) (*dto.BidResponse, error)
	SendDesicionByBid(req *dto.SendBidDesicionRequest) (*dto.BidResponse, error)
	SendFeedbackByBid(req *dto.SendBidFeedbackRequest) (*dto.BidResponse, error)
	RollbackBidById(bidId, username string, version int) (*dto.BidResponse, error)
	GetFeedbacksByTenderAndAuthor(tenderID, authorUsername, requesterUsername string, limit int, offset int) ([]*dto.BidFeedbackResponse, error)
	ConvertBidRequestToModel(bidReq dto.BidRequest) models.Bid
}

type bid struct {
	cfg config.Config
	log logger.ILogger
	db  storage.IStorage
}

func newBid(cfg config.Config, log logger.ILogger, db storage.IStorage) *bid {
	return &bid{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (b bid) GetBidsByFiltering(limit, offset int, username string) ([]dto.BidResponse, error) {
	const fn = "service.bid.GetBidsByFiltering"

	var responses []dto.BidResponse

	authorId, err := b.db.Employee().GetIdByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	bids, err := b.db.Bid().GetBidsByFiltering(limit, offset, authorId)
	if err != nil {
		return nil, err
	}

	for _, bid := range bids {
		responses = append(responses, dto.BidResponse{
			ID:   bid.ID,
			Name: bid.Name,
		})
	}
	return responses, nil
}

func (b bid) GetBidsByUsername(username string) ([]dto.BidResponse, error) {
	const fn = "service.bid.GetBidsByUsername"
	var responses []dto.BidResponse

	authorId, err := b.db.Employee().GetIdByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	bids, err := b.db.Bid().GetBidsByUsername(authorId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if len(bids) == 0 {
		return nil, customerrors.ErrBidNotFound
	}

	for _, bid := range bids {
		responses = append(responses, dto.BidResponse{
			ID:         bid.ID,
			Name:       bid.Name,
			Status:     bid.Status,
			AuthorType: bid.AuthorType,
			AuthorID:   bid.AuthorID,
			Version:    bid.Version,
			CreatedAt:  bid.CreatedAt,
		})
	}
	return responses, nil
}

func (b bid) CreateBid(bidRequest dto.BidRequest) (*dto.BidResponse, error) {
	const fn = "service.bid.CreateBid"

	bidModel := b.ConvertBidRequestToModel(bidRequest)
	createdBid, err := b.db.Bid().CreateBid(bidModel)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	bidModel.ID = createdBid.ID

	err = b.db.BidHistory().SaveBidHistory(bidModel)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &dto.BidResponse{
		ID:         createdBid.ID,
		Name:       createdBid.Name,
		Status:     createdBid.Status,
		AuthorType: createdBid.AuthorType,
		AuthorID:   createdBid.AuthorID,
		Version:    createdBid.Version,
		CreatedAt:  createdBid.CreatedAt,
	}, nil
}

func (b bid) GetBidListByTenderId(tenderId, username string) ([]dto.BidResponse, error) {
	const fn = "service.bid.GetBidListByTenderId"

	var responses []dto.BidResponse

	authorId, err := b.db.Employee().GetIdByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	bids, err := b.db.Bid().GetBidListByTenderId(tenderId, authorId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if len(bids) == 0 {
		return nil, customerrors.ErrBidNotFound
	}

	for _, bid := range bids {
		responses = append(responses, dto.BidResponse{
			ID:         bid.ID,
			Name:       bid.Name,
			Status:     bid.Status,
			AuthorType: bid.AuthorType,
			AuthorID:   bid.AuthorID,
			Version:    bid.Version,
			CreatedAt:  bid.CreatedAt,
		})
	}
	return responses, nil
}

func (b bid) GetBidStatus(bidId, username string) (status string, err error) {
	const fn = "service.bid.GetBidStatus"

	authorId, err := b.db.Employee().GetIdByUsername(username)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	status, err = b.db.Bid().GetBidStatus(bidId, authorId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return status, nil
}

func (b bid) ChangeBidStatus(bidReq *dto.BidStatusRequest) (*dto.BidResponse, error) {
	const fn = "service.bid.ChangeBidStatus"

	authorId, err := b.db.Employee().GetIdByUsername(bidReq.Username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	version, err := b.db.Bid().GetVersionBidById(bidReq.BidID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	version += 1

	bid := &models.Bid{
		ID:       uuid.MustParse(bidReq.BidID),
		Status:   string(bidReq.Status),
		Version:  version,
		AuthorID: authorId,
	}

	updatedBid, err := b.db.Bid().ChangeBidStatus(bid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	err = b.db.BidHistory().SaveBidHistory(*updatedBid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &dto.BidResponse{
		ID:         updatedBid.ID,
		Name:       updatedBid.Name,
		Status:     updatedBid.Status,
		AuthorType: updatedBid.AuthorType,
		AuthorID:   updatedBid.AuthorID,
		Version:    updatedBid.Version,
		CreatedAt:  updatedBid.CreatedAt,
	}, nil
}

func (b bid) ChangeBidById(bidReq *dto.ChangeBidRequest) (*dto.BidResponse, error) {
	const fn = "service.bid.ChangeBidById"

	authorId, err := b.db.Employee().GetIdByUsername(bidReq.Username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	version, err := b.db.Bid().GetVersionBidById(bidReq.BidID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	version += 1

	bid := &models.Bid{
		ID:          uuid.MustParse(bidReq.BidID),
		Name:        bidReq.Name,
		Description: bidReq.Description,
		Version:     version,
		AuthorID:    authorId,
	}

	updatedBid, err := b.db.Bid().ChangeBidById(bid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	err = b.db.BidHistory().SaveBidHistory(*updatedBid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &dto.BidResponse{
		ID:         updatedBid.ID,
		Name:       updatedBid.Name,
		Status:     updatedBid.Status,
		AuthorType: updatedBid.AuthorType,
		AuthorID:   updatedBid.AuthorID,
		Version:    updatedBid.Version,
		CreatedAt:  updatedBid.CreatedAt,
	}, nil
}

func (b bid) SendDesicionByBid(req *dto.SendBidDesicionRequest) (*dto.BidResponse, error) {
	const fn = "service.bid.SendDesicionByBid"

	authorId, err := b.db.Employee().GetIdByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	version, err := b.db.Bid().GetVersionBidById(req.BidID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	version += 1

	bid := &models.Bid{
		ID: uuid.MustParse(req.BidID),
		Decision: sql.NullString{
			String: string(req.Decision),
		},
		Version:  version,
		AuthorID: authorId,
	}

	updatedBid, err := b.db.Bid().SendDesicionByBid(bid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	err = b.db.BidHistory().SaveBidHistory(*updatedBid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &dto.BidResponse{
		ID:         updatedBid.ID,
		Name:       updatedBid.Name,
		Status:     updatedBid.Status,
		AuthorType: updatedBid.AuthorType,
		AuthorID:   updatedBid.AuthorID,
		Version:    updatedBid.Version,
		CreatedAt:  updatedBid.CreatedAt,
	}, nil
}

func (b bid) SendFeedbackByBid(req *dto.SendBidFeedbackRequest) (*dto.BidResponse, error) {
	const fn = "service.bid.SendFeedbackByBid"

	authorId, err := b.db.Employee().GetIdByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	bid := &models.Bid{
		ID:       uuid.MustParse(req.BidID),
		AuthorID: authorId,
	}

	updatedBid, err := b.db.BidFeedbacks().SendFeedbackByBid(bid, req.BidFeedback)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &dto.BidResponse{
		ID:         updatedBid.ID,
		Name:       updatedBid.Name,
		Status:     updatedBid.Status,
		AuthorType: updatedBid.AuthorType,
		AuthorID:   updatedBid.AuthorID,
		Version:    updatedBid.Version,
		CreatedAt:  updatedBid.CreatedAt,
	}, nil
}

func (b bid) RollbackBidById(bidId, username string, version int) (*dto.BidResponse, error) {
	const fn = "service.bid.RollbackBidById"

	bid, err := b.db.Bid().RollbackBidById(bidId, username, version)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &dto.BidResponse{
		ID:         bid.ID,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorID:   bid.AuthorID,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt,
	}, nil
}

func (b bid) GetFeedbacksByTenderAndAuthor(tenderID, authorUsername, requesterUsername string, limit int, offset int) ([]*dto.BidFeedbackResponse, error) {
	const fn = "service.bid.GetFeedbacksByTenderAndAuthor"

	authorId, err := b.db.Employee().GetIdByUsername(authorUsername)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	requesterId, err := b.db.Employee().GetIdByUsername(requesterUsername)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	feedbacks, err := b.db.BidFeedbacks().GetFeedbacksByTenderAndAuthor(tenderID, authorId, requesterId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var feedbackResponses []*dto.BidFeedbackResponse
	for _, feedback := range feedbacks {
		feedbackResponses = append(feedbackResponses, &dto.BidFeedbackResponse{
			ID:        feedback.ID,
			BidID:     feedback.BidID,
			Feedback:  feedback.Feedback,
			CreatedAt: feedback.CreatedAt,
			AuthorID:  feedback.AuthorID,
		})
	}

	return feedbackResponses, nil
}

func (b bid) ConvertBidRequestToModel(bidReq dto.BidRequest) models.Bid {
	return models.Bid{
		ID:          uuid.New(),
		Name:        bidReq.Name,
		Description: bidReq.Description,
		Status:      "Created",
		TenderID:    bidReq.TenderId,
		AuthorType:  string(bidReq.AuthorType),
		AuthorID:    bidReq.AuthorId,
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
