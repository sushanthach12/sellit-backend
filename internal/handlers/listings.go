package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
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

func List(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for listings")

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

		rows, err := db.Query(
			`
			SELECT id, title, description, price, city, status, created_at, updated_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
			`,
			limit, skip,
		)
		if err != nil || rows.Err() != nil {
			log.Printf("Get Listing: %v", err)
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
				log.Printf("rows.Scan: %v", err)
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
}
