package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	customerrors "git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/errors"
	"github.com/google/uuid"
)

type organizationRepo struct {
	db *sql.DB
}

func newOrganizationRepo(db *sql.DB) *organizationRepo {
	return &organizationRepo{
		db: db,
	}
}

func (t *organizationRepo) CheckID(id uuid.UUID) (bool, error) {
	const fn = "storage.organizationRepo.CheckID"

	var count int
	err := t.db.QueryRow("SELECT COUNT(id) FROM organization WHERE id = $1", id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	if count == 0 {
		return false, fmt.Errorf("%s: %w", fn, customerrors.ErrOrganizationNotFound)
	}

	return count > 0, nil
}

func (t *organizationRepo) GetIdByUserId(userId uuid.UUID) (uuid.UUID, error) {
	const fn = "storage.organizationRepo.GetIdByUserId"

	var id uuid.UUID
	err := t.db.QueryRow("SELECT organization_id FROM organization_responsible WHERE user_id = $1",
		userId).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("%s: %w", fn, customerrors.ErrOrganizationNotFound)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", fn, err)
	}

	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%s: %w", fn, customerrors.ErrOrganizationNotFound)
	}

	return id, nil
}
