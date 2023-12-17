package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/lastbyte32/wb-level-0/internal/model"
)

type getter interface {
	GetOrderByID(ctx context.Context, id string) (model.Order, error)
	GetOrderFromDB(ctx context.Context, id string) (model.Order, error)
}
type handler struct {
	service getter
}

func New(service getter) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) OrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	components := strings.Split(path, "/")

	id := components[len(components)-1]
	if uuid.Validate(id) != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrderByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "failed to marshal order to JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
