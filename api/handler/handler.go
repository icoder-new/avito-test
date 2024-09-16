package handler

import (
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/config"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/service"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/logger"
)

type Handler struct {
	cfg config.Config
	log logger.ILogger
	svc service.IService
}

func NewHandler(cfg config.Config, log logger.ILogger, svc service.IService) *Handler {
	return &Handler{
		cfg: cfg,
		log: log,
		svc: svc,
	}
}
