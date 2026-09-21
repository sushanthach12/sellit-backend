package listings

import (
	"context"
	"database/sql"
	"fmt"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) ListingRepository {
	return &repository{db: db}
}

func (r *repository) Count(ctx context.Context) (int, error) {
	var totalItems int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM listings`).Scan(&totalItems)
	if err != nil {
		return 0, fmt.Errorf("count listings: %w", err)
	}
	return totalItems, nil
}

func (r *repository) List(ctx context.Context, limit, skip int) ([]Listing, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, description, price, city, status, created_at, updated_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("query listings: %w", err)
	}
	defer rows.Close()

	listings := []Listing{}
	for rows.Next() {
		var l Listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.Status, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan listing row: %w", err)
		}
		listings = append(listings, l)
	}

	// rows.Err() must be checked AFTER the loop, not before — it only
	// reflects errors encountered during iteration, not from QueryContext
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listing rows: %w", err)
	}

	return listings, nil
}

func (r *repository) Create(ctx context.Context, l *Listing) error {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO listings (title, description, price, city, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, created_at
	`, l.Title, l.Description, l.Price, l.City, "active")

	if err := row.Scan(&l.ID, &l.Title, &l.CreatedAt); err != nil {
		return fmt.Errorf("insert listing: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		return 0, fmt.Errorf("delete listing: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}
	return affected, nil
}
