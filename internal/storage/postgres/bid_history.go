package postgres

import (
	"database/sql"
	"fmt"

	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/models"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/logger"
)

type bidHistoryRepo struct {
	db  *sql.DB
	log logger.ILogger
}

func newBidHistoryRepo(db *sql.DB) *bidHistoryRepo {
	return &bidHistoryRepo{
		db:  db,
		log: logger.NewLogger("local"),
	}
}

func (b *bidHistoryRepo) SaveBidHistory(bid models.Bid) error {
	const fn = "storage.bidRepo.SaveBidHistory"

	historyQuery := `INSERT INTO bid_history (bid_id, name, description, status, author_type, author_id, version, decision, created_at, updated_at)
                     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := b.db.Exec(historyQuery,
		bid.ID,
		bid.Name,
		bid.Description,
		bid.Status,
		bid.AuthorType,
		bid.AuthorID,
		bid.Version,
		bid.Decision,
		bid.CreatedAt,
		bid.UpdatedAt)

	if err != nil {
		return fmt.Errorf("%s: %v", fn, err)
	}
	return nil
}

func (b *bidRepo) GetVersionBidById(bidId string) (version int, err error) {
	const fn = "storage.bidRepo.GetVersionById"
	query := "SELECT version FROM bids WHERE id = $1"

	err = b.db.QueryRow(query, bidId).Scan(&version)

	if err != nil {
		return 0, fmt.Errorf("%s: %v", fn, err)
	}

	return version, nil
}
