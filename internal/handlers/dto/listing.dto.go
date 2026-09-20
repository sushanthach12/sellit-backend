package handlers

import "time"

// Only decode the required fields unnecessary fields are omitted
type CreateListingPayloadDto struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	City        string  `json:"city"`
}

type CreateListingResponseDto struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type GetListingResponseDto struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	City        string    `json:"city"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
