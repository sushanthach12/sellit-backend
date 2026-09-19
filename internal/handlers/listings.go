package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/sushanthach12/sellit-backend/internal/constants"
)

// make sure the key names are starting with the capitals, otherwise during the rows.Scan they wont be assigned, because they would be private
type listing struct {
	ID          string    `json:"id"` // while encoding the ID will be converted to the id (this is called struct tags)
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	City        string    `json:"city"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Constructor pattern for dependency handling for the handlers
// this improves code quality
/*
	* instead of having
	List(db)
	DeleteListing(db)

	* we do
	listingHandler := NewListingHandler()
	listingHandler.List
	listingHandler.DeleteListing

*/
type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{db: db, logger: logger}
}

func (lh *ListingHandler) List(w http.ResponseWriter, r *http.Request) {
	// request scoped context
	// using context helps in zombie query handling,
	// it helps if the client request is cancelled or timeout, this helps in closing the db queries as, without this the query will be continuously running
	ctx := r.Context()

	lh.logger.Info("Received request for listings")

	w.Header().Set("Content-Type", "application/json")

	queryParams := r.URL.Query()

	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil || page < 1 {
		page = 1 // default
	}

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit < 1 {
		limit = 10 // default
	}

	skip := (page - 1) * limit

	rows, err := lh.db.QueryContext(
		ctx,
		`
			SELECT id, title, description, price, city, status, created_at, updated_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
			`,
		limit, skip,
	)
	if err != nil || rows.Err() != nil {
		lh.logger.Error("Get Listing:", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	listings := []listing{}

	for rows.Next() {
		var l listing

		// the order in which the variables are mentioned here is required
		// because the Scan will be assigning in the way we select the columns in the SELECT query above, otherwise it throws the error
		err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.Status, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			lh.logger.Error("rows.Scan:", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		listings = append(listings, l)
	}

	resp := constants.NewPaginatedResponse(listings, constants.Pagination{
		Page:       page,
		PageSize:   limit,
		TotalItems: len(listings),
		TotalPages: len(listings) / limit, // ! FIX THIS BY TOTAL RECORDS query
	})

	_ = json.NewEncoder(w).Encode(resp)
}

func (lh *ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	lh.logger.Info("Received request for delete listing")

	w.Header().Set("Content-Type", "application/json")

	// Not required since id there is id they "Method not allowed" will thrown
	listingId := r.PathValue("id")
	if listingId == "" {
		slog.Error("Listing id is required")
		http.Error(w, "Listing id is required", http.StatusBadRequest)
		return
	}

	result, err := lh.db.ExecContext(
		ctx,
		`DELETE FROM listings WHERE id = $1;`,
		listingId,
	)
	if err != nil {
		// log.Printf("delete error: %v", err)
		lh.logger.Error("Failed to delete:", "listing_id", listingId, "err", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// below is not required, and we should not send any message like "Record not found" due to security concern
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		lh.logger.Error("rows affected", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Record Deleted Successfully"))
}
