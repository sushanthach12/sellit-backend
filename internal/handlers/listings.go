package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/sushanthach12/sellit-backend/internal/constants"
	dto "github.com/sushanthach12/sellit-backend/internal/handlers/dto"
	"github.com/sushanthach12/sellit-backend/internal/httpx"
	"github.com/sushanthach12/sellit-backend/internal/middleware"
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
	requestId := middleware.GetRequestIdFromContext(ctx)

	lh.logger.Info("Received request for listings", "requestId", requestId)

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

	// get total record count first so we can compute total pages
	var totalItems int
	err = lh.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM listings`).Scan(&totalItems)
	if err != nil {
		lh.logger.Error("Get Listing count:", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

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
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

	defer rows.Close()

	listings := []dto.GetListingResponseDto{}

	for rows.Next() {
		var l listing

		// the order in which the variables are mentioned here is required
		// because the Scan will be assigning in the way we select the columns in the SELECT query above, otherwise it throws the error
		err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.Status, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			lh.logger.Error("rows.Scan:", "error", err)
			httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
			return
		}

		listings = append(listings, dto.GetListingResponseDto{
			Id:          l.ID,
			Title:       l.Title,
			Description: l.Description,
			Price:       l.Price,
			City:        l.City,
			Status:      l.Status,
			CreatedAt:   l.CreatedAt,
		})
	}

	// compute total pages using actual total record count, rounding up
	totalPages := (totalItems + limit - 1) / limit

	response := constants.NewPaginatedResponse(listings, constants.Pagination{
		Page:       page,
		PageSize:   limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	})

	httpx.WriteJSON(w, http.StatusOK, response)
}

func (lh *ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	lh.logger.Info("Received request for delete listing")

	// Not required since id there is id they "Method not allowed" will thrown
	listingId := r.PathValue("id")
	if listingId == "" {
		slog.Error("Listing id is required")
		httpx.Error(w, http.StatusBadRequest, "Listing id is required", httpx.CodeInvalidId)
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
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

	// below is not required, and we should not send any message like "Record not found" due to security concern
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		lh.logger.Error("rows affected", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

	httpx.WriteJSON(w, http.StatusNoContent, constants.NewResponse("Record Deleted Successfully", nil))
}

func (lh *ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.GetRequestIdFromContext(ctx)

	lh.logger.Info("Received Request for listing create")

	// scoping/retrieve only the required fields from the request
	var payload dto.CreateListingPayloadDto
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		lh.logger.Error("Malformed Payload:", "error", err)
		httpx.Error(w, http.StatusBadGateway, "Invalid Payload", httpx.CodeMalformedJson)
		return
	}

	// DATA Validation
	if err := payload.Validate(); err != nil {
		var validationErr *constants.ValidationError
		errors.As(err, &validationErr) // just to the Field variable for the field return, because the err will not have it, so this get the parent or the first error in the error tree

		lh.logger.Error("Validation failed:", "error", err.Error())
		httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.CodeValidationFailed, validationErr.Field)
		return
	}

	// QueryRowContext for expected to return at-least one row after create
	row := lh.db.QueryRowContext(
		ctx,
		`
		INSERT INTO listings (title, "description", price, city, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, created_at
		`,
		payload.Title,
		payload.Description,
		payload.Price,
		payload.City,
		"active",
	)

	if err := row.Err(); err != nil {
		lh.logger.Error("Error Create Listing:", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

	var result dto.CreateListingResponseDto
	if err := row.Scan(&result.ID, &result.Title, &result.CreatedAt); err != nil {
		lh.logger.Error("Failed to insert:", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

	response := constants.NewResponse(result, nil)

	lh.logger.Info("Listing created", "request_id", requestId)

	httpx.WriteJSON(w, http.StatusCreated, response)
}
