package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/models"
	customerrors "git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/errors"
	"github.com/google/uuid"
)

type bidFeedbacksRepo struct {
	db *sql.DB
}

func newBidFeedbacksRepo(db *sql.DB) *bidFeedbacksRepo {
	return &bidFeedbacksRepo{
		db: db,
	}
}

func (b *bidFeedbacksRepo) SendFeedbackByBid(bid *models.Bid, bidFeedback string) (*models.Bid, error) {
	const fn = "storage.bindRepo.SendFeedbackByBid"

	feedbackQuery := `INSERT INTO bid_feedbacks (bid_id, author_id, feedback, created_at)
                      VALUES ($1, $2, $3, NOW())`
	_, err := b.db.Exec(feedbackQuery, bid.ID, bid.AuthorID, bidFeedback)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", fn, err)
	}

	bidQuery := `SELECT id, name, status, author_type, author_id, version, created_at
                 FROM bids
                 WHERE id = $1`
	err = b.db.QueryRow(bidQuery, bid.ID).Scan(
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

func (b *bidFeedbacksRepo) GetFeedbacksByTenderAndAuthor(tenderID string, authorId, requesterId uuid.UUID, limit int, offset int) ([]*models.BidFeedback, error) {
	const fn = "storage.bidFeedbacksRepo.GetFeedbacksByTenderAndAuthor"

	tenderQuery := `SELECT id FROM tenders WHERE id = $1 AND responsible_id = $2`

	_, err := b.db.Exec(tenderQuery, tenderID, requesterId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, customerrors.ErrTenderNotFound)
		}
		return nil, fmt.Errorf("%s: %v", fn, err)
	}

	feedbackQuery := `SELECT id, bid_id, author_id, feedback, created_at FROM bid_feedbacks WHERE author_id = $1 LIMIT $2 OFFSET $3`

	rows, err := b.db.Query(feedbackQuery, authorId, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("%s: %v", fn, err)
	}
	defer rows.Close()

	var feedbacks []*models.BidFeedback

	for rows.Next() {
		feedback := &models.BidFeedback{}
		if err := rows.Scan(
			&feedback.ID,
			&feedback.BidID,
			&feedback.AuthorID,
			&feedback.Feedback,
			&feedback.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %v", fn, err)
		}
		feedbacks = append(feedbacks, feedback)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %v", fn, err)
	}

	return feedbacks, nil
}
