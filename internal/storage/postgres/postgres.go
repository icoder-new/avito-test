package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/config"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/storage"

	_ "github.com/lib/pq"
)

type store struct {
	db                *sql.DB
	repoTender        *tenderRepo
	bidRepo           *bidRepo
	repoOrganization  *organizationRepo
	repoEmployee      *employeeRepo
	bidHistoryRepo    *bidHistoryRepo
	bidFeedbacksRepo  *bidFeedbacksRepo
	repoTenderHistory *tenderHistoryRepo
}

func newStorage(db *sql.DB) *store {
	return &store{
		db:                db,
		repoTender:        newTenderRepo(db),
		bidRepo:           newBidRepo(db),
		repoOrganization:  newOrganizationRepo(db),
		repoEmployee:      newEmployeeRepo(db),
		bidHistoryRepo:    newBidHistoryRepo(db),
		bidFeedbacksRepo:  newBidFeedbacksRepo(db),
		repoTenderHistory: newTenderHistoryRepo(db),
	}
}

func NewStorage(ctx context.Context, cfg config.Config) (storage.IStorage, error) {
	const fn = "storage.postgres.NewStorage"

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return newStorage(db), nil
}

func (s *store) CloseDB() {
	s.db.Close()
}

func (s *store) Tender() storage.ITender {
	return s.repoTender
}

func (s *store) Bid() storage.IBid {
	return s.bidRepo
}

func (s *store) BidHistory() storage.IBidHistory {
	return s.bidHistoryRepo
}

func (s *store) Organization() storage.IOrganization {
	return s.repoOrganization
}

func (s *store) Employee() storage.IEmployee {
	return s.repoEmployee
}

func (s *store) BidFeedbacks() storage.IBidFeedbacks {
	return s.bidFeedbacksRepo
}

func (s *store) TenderHistory() storage.ITenderHistory {
	return s.repoTenderHistory
}
