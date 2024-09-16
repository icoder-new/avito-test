package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/models"
	customerrors "git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/errors"
	"github.com/google/uuid"
)

type bidRepo struct {
	db *sql.DB
}

func newBidRepo(db *sql.DB) *bidRepo {
	return &bidRepo{
		db: db,
	}
}

func (b *bidRepo) GetBidsByFiltering(limit, offset int, authorId uuid.UUID) ([]models.Bid, error) {
	const fn = "storage.bindRepo.GetBidsByFiltering"

	query := "SELECT id, name, status, author_type, author_id, version, created_at FROM bids WHERE author_id = $1 LIMIT $2 OFFSET $3"
	rows, err := b.db.Query(query, authorId, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, customerrors.ErrBidNotFound)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		err = rows.Scan(
			&bid.ID, &bid.Name, &bid.Status, &bid.AuthorType,
			&bid.AuthorID, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		bids = append(bids, bid)
	}
	return bids, nil
}

func (b *bidRepo) GetBidsByUsername(authorId uuid.UUID) ([]models.Bid, error) {
	const fn = "storage.bindRepo.GetBidsByUsername"

	query := "SELECT id, name, status, author_type, author_id, version, created_at FROM bids WHERE author_id = $1"
	rows, err := b.db.Query(query, authorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, customerrors.ErrBidNotFound)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		err = rows.Scan(
			&bid.ID, &bid.Name, &bid.Status, &bid.AuthorType,
			&bid.AuthorID, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		bids = append(bids, bid)
	}
	return bids, nil
}

func (b *bidRepo) CreateBid(bid models.Bid) (*models.Bid, error) {
	const fn = "storage.bidRepo.CreateBid"

	query := `INSERT INTO bids(name, description, status, tender_id, author_type, author_id, version, created_at, updated_at)
              VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
              RETURNING id, name, status, author_type, author_id, version, created_at, updated_at`

	var createdBid models.Bid
	err := b.db.QueryRow(query, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.AuthorType, bid.AuthorID, bid.Version, bid.CreatedAt, bid.UpdatedAt).Scan(
		&createdBid.ID,
		&createdBid.Name,
		&createdBid.Status,
		&createdBid.AuthorType,
		&createdBid.AuthorID,
		&createdBid.Version,
		&createdBid.CreatedAt,
		&createdBid.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", fn, err)
	}
	return &createdBid, nil
}

func (b *bidRepo) GetBidListByTenderId(tenderId string, authorId uuid.UUID) ([]models.Bid, error) {
	const fn = "storage.bindRepo.GetBidListByTenderId"

	query := "SELECT id, name, status, author_type, author_id, version, created_at FROM bids WHERE tender_id = $1 AND author_id = $2"
	rows, err := b.db.Query(query, tenderId, authorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, customerrors.ErrBidNotFound)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		err = rows.Scan(
			&bid.ID, &bid.Name, &bid.Status, &bid.AuthorType,
			&bid.AuthorID, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		bids = append(bids, bid)
	}
	return bids, nil
}

func (b *bidRepo) GetBidStatus(bidId string, authorId uuid.UUID) (status string, err error) {
	const fn = "storage.bindRepo.GetBidStatus"

	query := "SELECT status FROM bids WHERE id = $1 AND author_id = $2"
	row := b.db.QueryRow(query, bidId, authorId)

	err = row.Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", fn, customerrors.ErrBidNotFound)
		}
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return status, nil
}

func (b *bidRepo) ChangeBidStatus(bid *models.Bid) (*models.Bid, error) {
	const fn = "storage.bindRepo.ChangeBidStatus"

	query := "UPDATE bids SET status = $1, version = $2, updated_at = NOW() WHERE id = $3 AND author_id = $4 RETURNING id, name, status, author_type, author_id, version, created_at"

	err := b.db.QueryRow(query, bid.Status, bid.Version, bid.ID, bid.AuthorID).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Status,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Version,
		&bid.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %v", fn, err)
	}

	return bid, nil
}

func (b *bidRepo) ChangeBidById(bid *models.Bid) (*models.Bid, error) {
	const fn = "storage.bindRepo.ChangeBidByBidId"

	query := "UPDATE bids SET name = $1, description = $2, version = $3 WHERE id = $4 AND author_id = $5 RETURNING id, name, status, author_type, author_id, version, created_at"

	err := b.db.QueryRow(query, bid.Name, bid.Description, bid.Version, bid.ID, bid.AuthorID).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Status,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Version,
		&bid.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %v", fn, err)
	}
	return bid, nil
}

func (b *bidRepo) SendDesicionByBid(bid *models.Bid) (*models.Bid, error) {
	const fn = "storage.bindRepo.SendDesicionByBid"

	query := "UPDATE bids SET decision = $1, version = $2, updated_at = NOW() WHERE id = $3 AND author_id = $4 RETURNING id, name, status, author_type, author_id, version, decision, created_at"

	err := b.db.QueryRow(query, bid.Decision, bid.Version, bid.ID, bid.AuthorID).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Status,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Version,
		&bid.Decision,
		&bid.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %v", fn, err)
	}
	return bid, nil
}

func (b *bidRepo) RollbackBidById(bidId, username string, version int) (*models.Bid, error) {
	const fn = "storage.bidRepo.RollbackBidById"

	query := `SELECT bh.id, bh.name, bh.description, bh.status, bh.author_type, bh.author_id, b.version + 1 AS version, bh.decision, bh.created_at FROM bid_history bh JOIN bids b ON bh.bid_id = b.id WHERE bh.bid_id = $1 AND bh.version = $2`

	var bid models.Bid
	err := b.db.QueryRow(query, bidId, version).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Version,
		&bid.Decision,
		&bid.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, customerrors.ErrBidNotFound)
		}
		return nil, fmt.Errorf("%s: %v", fn, err)
	}

	newVersion := bid.Version + 1
	updateQuery := `UPDATE bids 
					SET name = $1, description = $2, status = $3, author_type = $4, author_id = $5, version = $6, decision = $7, created_at = $8, updated_at = NOW()
					WHERE id = $9 AND author_id = (SELECT id FROM employee WHERE username = $10)`
	_, err = b.db.Exec(updateQuery, bid.Name, bid.Description, bid.Status, bid.AuthorType, bid.AuthorID, newVersion, bid.Decision, bid.CreatedAt, bidId, username)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", fn, err)
	}

	bid.Version = newVersion

	return &bid, nil
}
