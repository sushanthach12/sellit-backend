package listings

import (
	"time"

	"github.com/sushanthach12/sellit-backend/internal/constants"
	"github.com/sushanthach12/sellit-backend/internal/helpers"
)

// Only decode the required fields unnecessary fields are omitted
type CreateListingPayloadDto struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	City        string  `json:"city"`
}

func (payload *CreateListingPayloadDto) Validate() error {
	if helpers.CheckIfStringEmpty(payload.Title) {
		return &constants.ValidationError{
			Field:   "title",
			Message: "must not be empty",
		}
	}

	if helpers.CheckIfStringEmpty(payload.Description) {
		return &constants.ValidationError{
			Field:   "description",
			Message: "must not be empty",
		}
	}

	// if helpers.CheckStringLen(payload.Description, 1, 10000) {
	// 	return &constants.ValidationError{
	// 		Field:   "description",
	// 		Message: "must be within 1 - 5000",
	// 	}
	// }

	if helpers.CheckIfStringEmpty(payload.City) {
		return &constants.ValidationError{
			Field:   "city",
			Message: "must not be empty",
		}
	}

	if !helpers.CheckIfValidNumber(payload.Price, true) {
		return &constants.ValidationError{
			Field:   "price",
			Message: "should be greater than 0",
		}
	}

	return nil
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
