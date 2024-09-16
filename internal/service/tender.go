package service

import (
	"errors"
	"fmt"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/config"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/dto"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/models"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/storage"
	customerrors "git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/errors"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/logger"
	"github.com/google/uuid"
	"time"
)

type ITender interface {
	GetAllTenders(tender dto.TenderGetRequest) ([]dto.TenderResponse, error)
	CreateTender(tender dto.TenderRequest) (dto.TenderResponse, error)
	GetMyTenders(tender dto.TenderMyGetRequest) ([]dto.TenderResponse, error)
	GetTenderStatusById(tender dto.TenderStatusRequest) (string, error)
	ChangeTenderStatusById(tenderStatus dto.TenderStatusRequest) (dto.TenderResponse, error)
	EditTenderById(tenderEdit dto.TenderEditRequest) (dto.TenderResponse, error)
	RollbackTenderById(rbTender dto.TenderRollbackRequest) (dto.TenderResponse, error)
}

type tender struct {
	cfg config.Config
	log logger.ILogger
	db  storage.IStorage
}

func newTender(cfg config.Config, log logger.ILogger, db storage.IStorage) *tender {
	return &tender{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (t tender) GetAllTenders(tender dto.TenderGetRequest) ([]dto.TenderResponse, error) {
	const fn = "service.tender.GetAllTenders"

	var responses []dto.TenderResponse

	tenders, err := t.db.Tender().GetAllTenders(tender.Limit, tender.Offset, string(tender.ServiceType))
	if err != nil {
		t.log.Error(fn, err)
		return nil, err
	}

	for _, tender := range tenders {
		responses = append(responses, dto.TenderResponse{
			ID:          tender.ID,
			Name:        tender.Name,
			Description: tender.Description,
			Status:      tender.Status,
			ServiceType: tender.ServiceType,
			Version:     tender.Version,
			CreatedAt:   tender.CreatedAt,
		})
	}

	if len(responses) == 0 {
		return responses, fmt.Errorf("%s: %w", fn, customerrors.ErrTenderNotFound)
	}

	return responses, nil
}

func (t tender) CreateTender(tender dto.TenderRequest) (dto.TenderResponse, error) {
	const fn = "service.tender.CreateTender"

	userId, err := t.db.Employee().GetIdByUsername(tender.CreatorUsername)
	if err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	if ok, err := t.db.Organization().CheckID(tender.OrganizationID); !ok || err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	if ok, err := t.db.Employee().IsOwnerOrganization(userId, tender.OrganizationID); !ok || err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	mTender := models.Tender{
		ID:             uuid.New(),
		Name:           tender.Name,
		Description:    tender.Description,
		Status:         "Created",
		ServiceType:    string(tender.ServiceType),
		OrganizationID: tender.OrganizationID,
		ResponsibleID:  userId,
		Version:        1,
		CreatedAt:      time.Now(),
	}

	if err := t.db.TenderHistory().CreateTenderHistory(uuid.New(), mTender); err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	response, err := t.db.Tender().CreateTender(mTender)
	if err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	return dto.TenderResponse{
		ID:          response.ID,
		Name:        response.Name,
		Description: response.Description,
		Status:      response.Status,
		ServiceType: response.ServiceType,
		Version:     response.Version,
		CreatedAt:   response.CreatedAt,
	}, nil
}

func (t tender) GetMyTenders(tender dto.TenderMyGetRequest) ([]dto.TenderResponse, error) {
	const fn = "service.tender.GetMyTenders"

	userId, err := t.db.Employee().GetIdByUsername(tender.Username)
	if err != nil {
		t.log.Error(fn, err)
		return nil, err
	}

	organizationId, err := t.db.Organization().GetIdByUserId(userId)
	if err != nil {
		t.log.Error(fn, err)
		return nil, err
	}

	tenders, err := t.db.Tender().GetMyTenders(tender.Limit, tender.Offset, userId, organizationId)
	if err != nil {
		t.log.Error(fn, err)
		return nil, err
	}

	if len(tenders) == 0 {
		t.log.Error(fn, fmt.Errorf("%s: %w", fn, customerrors.ErrTenderNotFound))
		return nil, fmt.Errorf("%s: %w", fn, customerrors.ErrTenderNotFound)
	}

	var responses []dto.TenderResponse
	for _, tender := range tenders {
		responses = append(responses, dto.TenderResponse{
			ID:          tender.ID,
			Name:        tender.Name,
			Description: tender.Description,
			Status:      tender.Status,
			ServiceType: tender.ServiceType,
			Version:     tender.Version,
			CreatedAt:   tender.CreatedAt,
		})
	}

	return responses, nil
}

func (t tender) GetTenderStatusById(tender dto.TenderStatusRequest) (string, error) {
	const fn = "service.tender.GetTenderStatusById"

	userId, err := t.db.Employee().GetIdByUsername(tender.Username)
	if err != nil {
		t.log.Error(fn, err)
		return "", err
	}

	status, err := t.db.Tender().GetTenderStatusById(tender.TenderId, userId.String())
	if err != nil {
		t.log.Error(fn, err)
		return "", err
	}

	return status, nil
}

func (t tender) ChangeTenderStatusById(tenderStatus dto.TenderStatusRequest) (dto.TenderResponse, error) {
	const fn = "service.tender.ChangeTenderStatusById"

	userId, err := t.db.Employee().GetIdByUsername(tenderStatus.Username)
	if err != nil && !errors.Is(err, customerrors.ErrEmployeeNotFound) {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	if ok, err := t.db.Tender().CheckID(tenderStatus.TenderId); !ok || err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	if ok, err := t.db.Employee().IsTenderOwner(userId, tenderStatus.TenderId); !ok || err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	response, err := t.db.Tender().GetTenderById(tenderStatus.TenderId)
	if err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	response.Status = string(tenderStatus.Status)
	response.Version++
	response.UpdatedAt = time.Now()

	if err := t.db.TenderHistory().CreateTenderHistory(uuid.New(), response); err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	if err := t.db.Tender().UpdateTender(response); err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	return dto.TenderResponse{
		ID:          response.ID,
		Name:        response.Name,
		Description: response.Description,
		Status:      response.Status,
		ServiceType: response.ServiceType,
		Version:     response.Version,
		CreatedAt:   response.CreatedAt,
	}, nil
}

func (t tender) EditTenderById(tenderEdit dto.TenderEditRequest) (dto.TenderResponse, error) {
	const fn = "service.tender.EditTenderById"

	userId, err := t.db.Employee().GetIdByUsername(tenderEdit.Username)
	if err != nil && !errors.Is(err, customerrors.ErrEmployeeNotFound) {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	if ok, err := t.db.Tender().CheckID(tenderEdit.TenderId); !ok || err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	if ok, err := t.db.Employee().IsTenderOwner(userId, tenderEdit.TenderId); !ok || err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	response, err := t.db.Tender().GetTenderById(tenderEdit.TenderId)
	if err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	if tenderEdit.Name != "" && tenderEdit.Name != " " {
		response.Name = tenderEdit.Name
	}

	if tenderEdit.Description != "" && tenderEdit.Description != " " {
		response.Description = tenderEdit.Description
	}

	response.ServiceType = string(tenderEdit.ServiceType)

	response.Version++
	response.UpdatedAt = time.Now()

	if err := t.db.TenderHistory().CreateTenderHistory(uuid.New(), response); err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	if err := t.db.Tender().UpdateTender(response); err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	return dto.TenderResponse{
		ID:          response.ID,
		Name:        response.Name,
		Description: response.Description,
		Status:      response.Status,
		ServiceType: response.ServiceType,
		Version:     response.Version,
		CreatedAt:   response.CreatedAt,
	}, nil
}

func (t tender) RollbackTenderById(rbTender dto.TenderRollbackRequest) (dto.TenderResponse, error) {
	const fn = "service.tender.RollbackTenderById"

	userId, err := t.db.Employee().GetIdByUsername(rbTender.Username)
	if err != nil && !errors.Is(err, customerrors.ErrEmployeeNotFound) {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	if ok, err := t.db.Tender().CheckID(rbTender.TenderId); !ok || err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	if ok, err := t.db.Employee().IsTenderOwner(userId, rbTender.TenderId); !ok || err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	organizationId, err := t.db.Organization().GetIdByUserId(userId)
	if err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	oldTender, err := t.db.Tender().GetTenderById(rbTender.TenderId)
	if err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	response, err := t.db.TenderHistory().GetTender(rbTender.Version, userId, organizationId)
	if err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	oldTender.Name = response.Name
	oldTender.Description = response.Description
	oldTender.Status = response.Status
	oldTender.ServiceType = response.ServiceType
	oldTender.Version++
	oldTender.CreatedAt = response.CreatedAt
	oldTender.UpdatedAt = time.Now()

	if err := t.db.TenderHistory().CreateTenderHistory(uuid.New(), oldTender); err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	if err := t.db.Tender().UpdateTender(oldTender); err != nil {
		t.log.Error(fn, err)
		return dto.TenderResponse{}, err
	}

	return dto.TenderResponse{
		ID:          oldTender.ID,
		Name:        oldTender.Name,
		Description: oldTender.Description,
		Status:      oldTender.Status,
		ServiceType: oldTender.ServiceType,
		Version:     oldTender.Version,
		CreatedAt:   oldTender.CreatedAt,
	}, nil
}
