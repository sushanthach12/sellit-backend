package listings

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/sushanthach12/sellit-backend/internal/constants"
	"github.com/sushanthach12/sellit-backend/internal/helpers"
	"github.com/sushanthach12/sellit-backend/internal/httpx"
	"github.com/sushanthach12/sellit-backend/internal/middleware"
)

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
type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(svc Service, logger *slog.Logger) *Handler {
	return &Handler{service: svc, logger: logger}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	params := helpers.ParsePaginationParams(r) // your existing page/limit parsing

	response, err := h.service.List(r.Context(), params.Page, params.Limit)
	if err != nil {
		h.logger.Error("list listings failed", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, response)
}

func (lh *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.GetRequestIdFromContext(ctx)

	lh.logger.Info("Received Request for listing create")

	var payload CreateListingPayloadDto
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		lh.logger.Error("Malformed Payload:", "error", err)
		httpx.Error(w, http.StatusBadRequest, "Invalid Payload", httpx.CodeMalformedJson)
		return
	}

	result, err := lh.service.Create(ctx, payload)
	if err != nil {
		if validationErr, ok := errors.AsType[*constants.ValidationError](err); ok {
			lh.logger.Error("Validation failed:", "error", err.Error())
			httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.CodeValidationFailed, validationErr.Field)
			return
		}

		lh.logger.Error("Failed to create listing:", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

	response := constants.NewResponse(result, nil)
	lh.logger.Info("Listing created", "request_id", requestId)
	httpx.WriteJSON(w, http.StatusCreated, response)
}

func (lh *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	lh.logger.Info("Received request for delete listing")

	listingId := r.PathValue("id")
	if listingId == "" {
		lh.logger.Error("Listing id is required")
		httpx.Error(w, http.StatusBadRequest, "Listing id is required", httpx.CodeInvalidId)
		return
	}

	err := lh.service.Delete(ctx, listingId)
	if err != nil {
		if errors.Is(err, ErrListingNotFound) {
			// keep response generic — don't confirm/deny existence, per your own security note
			httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
			return
		}

		lh.logger.Error("Failed to delete:", "listing_id", listingId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong!", httpx.CodeInternalError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
