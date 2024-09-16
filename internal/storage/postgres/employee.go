package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	customerrors "git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/errors"
	"github.com/google/uuid"
)

type employeeRepo struct {
	db *sql.DB
}

func newEmployeeRepo(db *sql.DB) *employeeRepo {
	return &employeeRepo{
		db: db,
	}
}

func (t *employeeRepo) GetIdByUsername(username string) (uuid.UUID, error) {
	const fn = "storage.employeeRepo.GetIdByUsername"

	var id uuid.UUID
	err := t.db.QueryRow("SELECT id FROM employee WHERE username = $1", username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("%s: %w", fn, customerrors.ErrEmployeeNotFound)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", fn, err)
	}

	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%s: %w", fn, customerrors.ErrEmployeeNotFound)
	}

	return id, nil
}

func (t *employeeRepo) IsOwnerOrganization(id, organizationId uuid.UUID) (bool, error) {
	const fn = "storage.employeeRepo.IsOwnerOrganization"

	var count int

	query := `SELECT COUNT(id) FROM organization_responsible WHERE user_id = $1 AND organization_id = $2`
	err := t.db.QueryRow(query, id, organizationId).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", fn, customerrors.ErrNotEnoughRights)
		}

		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return count > 0, nil
}

func (t *employeeRepo) IsTenderOwner(id, tenderId uuid.UUID) (bool, error) {
	const fn = "storage.employeeRepo.IsTenderOwner"

	var count int

	query := `SELECT COUNT(*) FROM tenders WHERE responsible_id = $1 AND id = $2`
	err := t.db.QueryRow(query, id, tenderId).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", fn, customerrors.ErrNotEnoughRights)
		}
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return count > 0, nil
}
