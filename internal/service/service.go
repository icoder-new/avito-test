package service

import (
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/config"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/storage"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/logger"
)

type IService interface {
	Tender() ITender
	Bid() IBid
}

type service struct {
	tender ITender
	bid    IBid
}

func (s service) Tender() ITender {
	return s.tender
}

func (s service) Bid() IBid {
	return s.bid
}

func NewService(cfg config.Config, log logger.ILogger, db storage.IStorage) IService {
	return &service{
		tender: newTender(cfg, log, db),
		bid:    newBid(cfg, log, db),
	}
}
