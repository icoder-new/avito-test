package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/models"
	customerrors "git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/errors"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/logger"
	"github.com/google/uuid"
)

type tenderHistoryRepo struct {
	db  *sql.DB
	log logger.ILogger
}

func newTenderHistoryRepo(db *sql.DB) *tenderHistoryRepo {
	return &tenderHistoryRepo{
		db:  db,
		log: logger.NewLogger("local"),
	}
}

func (t *tenderHistoryRepo) CreateTenderHistory(id uuid.UUID, tender models.Tender) error {
	const fn = "storage.tenderHistoryRepo.CreateTenderHistory"

	query := `INSERT INTO tender_history
    		(id, name, description, service_type, status, organization_id, responsible_id, version, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := t.db.Exec(
		query,
		id,
		tender.Name,
		tender.Description,
		tender.ServiceType,
		tender.Status,
		tender.OrganizationID,
		tender.ResponsibleID,
		tender.Version,
		tender.CreatedAt,
		tender.UpdatedAt,
	)
	if err != nil {
		t.log.Error(fn, err)
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (t *tenderHistoryRepo) GetVersionById(tenderId uuid.UUID) (int, error) {
	const fn = "storage.tenderHistoryRepo.GetVersionById"

	var version int

	err := t.db.QueryRow("SELECT version FROM tender_history WHERE id = $1", tenderId).Scan(&version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", fn, customerrors.ErrTenderNotFound)
		}

		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return version, nil
}

func (t *tenderHistoryRepo) GetTender(version int, userId, organizationId uuid.UUID) (models.Tender, error) {
	const fn = "storage.tenderHistoryRepo.GetTender"

	var tender models.Tender

	err := t.db.QueryRow("SELECT * FROM tender_history WHERE version = $1 and responsible_id = $2 and organization_id = $3",
		version, userId, organizationId).
		Scan(
			&tender.ID,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.OrganizationID,
			&tender.ResponsibleID,
			&tender.Version,
			&tender.CreatedAt,
			&tender.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Tender{}, fmt.Errorf("%s: %w", fn, customerrors.ErrTenderNotFound)
		}
		return models.Tender{}, fmt.Errorf("%s: %w", fn, err)
	}

	return tender, nil
}
