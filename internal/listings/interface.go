package listings

import (
	"context"

	"github.com/sushanthach12/sellit-backend/internal/constants"
)

type Service interface {
	List(ctx context.Context, page, limit int) (constants.Response[GetListingResponseDto], error)
	Create(ctx context.Context, payload CreateListingPayloadDto) (CreateListingResponseDto, error)
	Delete(ctx context.Context, id string) error
}

type ListingRepository interface {
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, skip int) ([]Listing, error)
	Create(ctx context.Context, l *Listing) error
	Delete(ctx context.Context, id string) (int64, error)
}
