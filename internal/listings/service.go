package listings

import (
	"context"
	"errors"
	"fmt"

	"github.com/sushanthach12/sellit-backend/internal/constants"
)

// Sentinel errors returned by the listing service.
var (
	ErrListingNotFound = errors.New("listing not found")
)

type service struct {
	repository ListingRepository
}

func NewService(repo ListingRepository) Service {
	return &service{repository: repo}
}

func (s *service) List(ctx context.Context, page, limit int) (constants.Response[GetListingResponseDto], error) {
	skip := (page - 1) * limit

	totalItems, err := s.repository.Count(ctx)
	if err != nil {
		return constants.Response[GetListingResponseDto]{}, fmt.Errorf("service: get listing count: %w", err)
	}

	listings, err := s.repository.List(ctx, limit, skip)
	if err != nil {
		return constants.Response[GetListingResponseDto]{}, fmt.Errorf("service: list listings: %w", err)
	}

	items := make([]GetListingResponseDto, 0, len(listings))
	for _, l := range listings {
		items = append(items, GetListingResponseDto{
			Id:          l.ID,
			Title:       l.Title,
			Description: l.Description,
			Price:       l.Price,
			City:        l.City,
			Status:      l.Status,
			CreatedAt:   l.CreatedAt,
		})
	}

	totalPages := (totalItems + limit - 1) / limit

	return constants.NewPaginatedResponse(items, constants.Pagination{
		Page:       page,
		PageSize:   limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}), nil
}

func (s *service) Create(ctx context.Context, payload CreateListingPayloadDto) (CreateListingResponseDto, error) {
	if err := payload.Validate(); err != nil {
		return CreateListingResponseDto{}, err // validation error passed through as-is
	}

	l := &Listing{
		Title:       payload.Title,
		Description: payload.Description,
		Price:       payload.Price,
		City:        payload.City,
	}

	if err := s.repository.Create(ctx, l); err != nil {
		return CreateListingResponseDto{}, fmt.Errorf("service: create listing: %w", err)
	}

	return CreateListingResponseDto{
		ID:        l.ID,
		Title:     l.Title,
		CreatedAt: l.CreatedAt,
	}, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	affected, err := s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service: delete listing: %w", err)
	}

	if affected == 0 {
		return ErrListingNotFound
	}

	return nil
}
