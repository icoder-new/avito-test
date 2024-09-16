package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/models"
	customerrors "git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/errors"
	"github.com/google/uuid"
	"time"
)

type tenderRepo struct {
	db *sql.DB
}

func newTenderRepo(db *sql.DB) *tenderRepo {
	return &tenderRepo{
		db: db,
	}
}

func (t *tenderRepo) CreateTender(tender models.Tender) (models.Tender, error) {
	const fn = "storage.tenderRepo.CreateTender"

	var createdTender models.Tender

	query := `INSERT INTO tenders (id, name, description, status, service_type, organization_id, responsible_id, version, created_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
					RETURNING id, name, description, status, service_type, version, created_at`

	row := t.db.QueryRow(query,
		tender.ID, tender.Name, tender.Description, tender.Status, tender.ServiceType,
		tender.OrganizationID, tender.ResponsibleID, tender.Version, tender.CreatedAt.Format(time.RFC3339))
	if row == nil {
		return models.Tender{}, fmt.Errorf("%s: %w", fn, customerrors.ErrTenderNotFound)
	}

	err := row.Scan(&createdTender.ID, &createdTender.Name, &createdTender.Description,
		&createdTender.Status, &createdTender.ServiceType, &createdTender.Version, &createdTender.CreatedAt)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", fn, err)
	}

	return createdTender, nil
}

func (t *tenderRepo) GetAllTenders(limit, offset int, serviceType string) ([]models.Tender, error) {
	const fn = "storage.tenderRepo.GetAllTenders"

	var (
		query string
		rows  *sql.Rows
		err   error
	)

	if serviceType == "" || serviceType == " " {
		query = `SELECT id, name, description, status, service_type, version, created_at
				FROM tenders where status = 'PUBLISHED' LIMIT $1 OFFSET $2`
		rows, err = t.db.Query(query, limit, offset)
	} else {
		query = `SELECT id, name, description, status, service_type, version, created_at
				FROM tenders where status = 'PUBLISHED' and service_type = $1 LIMIT $2 OFFSET $3`
		rows, err = t.db.Query(query, serviceType, limit, offset)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, customerrors.ErrTenderNotFound)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var tenders []models.Tender
	for rows.Next() {
		var tender models.Tender
		err = rows.Scan(&tender.ID, &tender.Name, &tender.Description,
			&tender.Status, &tender.ServiceType, &tender.Version, &tender.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (t *tenderRepo) GetMyTenders(limit, offset int, userId, organizationId uuid.UUID) ([]models.Tender, error) {
	const fn = "storage.tenderRepo.GetMyTenders"

	query := `SELECT id, name, description, status, service_type, version, created_at
				FROM tenders where responsible_id = $1 and organization_id = $2 LIMIT $3 OFFSET $4`
	rows, err := t.db.Query(query, userId, organizationId, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, customerrors.ErrTenderNotFound)
		}

		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var tenders []models.Tender
	for rows.Next() {
		var tender models.Tender
		err = rows.Scan(&tender.ID, &tender.Name, &tender.Description,
			&tender.Status, &tender.ServiceType, &tender.Version, &tender.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (t *tenderRepo) GetTenderStatusById(tenderId uuid.UUID, userId string) (string, error) {
	const fn = "storage.tenderRepo.GetTenderStatusById"

	var (
		query  string
		rows   *sql.Rows
		err    error
		status string
	)

	if userId == uuid.Nil.String() {
		query = `SELECT status
				FROM tenders where id = $1`
		rows, err = t.db.Query(query, tenderId)
	} else {
		query = `SELECT status
				FROM tenders where id = $1 and responsible_id = $2 and organization_id = (SELECT organization_id FROM organization_responsible WHERE user_id = $3)`
		rows, err = t.db.Query(query, tenderId, userId, userId)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", fn, customerrors.ErrTenderNotFound)
		}
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	for rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			return "", fmt.Errorf("%s: %w", fn, err)
		}
	}

	return status, nil
}

func (t *tenderRepo) CheckID(tenderId uuid.UUID) (bool, error) {
	const fn = "storage.tenderRepo.CheckID"

	var count int

	query := `SELECT count(*) FROM tenders where id = $1`
	err := t.db.QueryRow(query, tenderId).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", fn, customerrors.ErrTenderNotFound)
		}
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return count > 0, nil
}

func (t *tenderRepo) GetTenderById(tenderId uuid.UUID) (models.Tender, error) {
	const fn = "storage.tenderRepo.GetTenderById"

	var tender models.Tender
	query := `SELECT id, name, description, service_type, status,
       			organization_id, responsible_id, version, created_at, updated_at
				FROM tenders where id = $1`
	err := t.db.QueryRow(query, tenderId).Scan(
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
			return tender, fmt.Errorf("%s: %w", fn, customerrors.ErrTenderNotFound)
		}
		return tender, fmt.Errorf("%s: %w", fn, err)
	}

	return tender, nil
}

func (t *tenderRepo) UpdateTender(tender models.Tender) error {
	const fn = "storage.tenderRepo.UpdateTender"

	query := `UPDATE tenders
				SET name = $1, description = $2, service_type = $3, status = $4, version = $5
				WHERE id = $6`
	_, err := t.db.Exec(
		query,
		tender.Name,
		tender.Description,
		tender.ServiceType,
		tender.Status,
		tender.Version,
		tender.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
