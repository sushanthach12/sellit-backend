package listings

import "time"

// make sure the key names are starting with the capitals, otherwise during the rows.Scan they wont be assigned, because they would be private
type Listing struct {
	ID          string    `json:"id"` // while encoding the ID will be converted to the id (this is called struct tags)
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	City        string    `json:"city"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
