package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// getOrder godoc
// @Summary      Получить ответ.
// @Param 		  code   path    int     true        "Код ответа"
// @Success      200  "success"
// @Failure      400  "error"
// @Failure      401  "error"
// @Failure      403  "error"
// @Failure      404  "error"
// @Failure      429  "error"
// @Failure      500  "error"
// @Failure      503  "error"
// @Router       /get/{code} [get]
func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.Atoi(chi.URLParam(r, "code"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	switch orderID {
	case 400:
		w.WriteHeader(http.StatusBadRequest)
	case 401:
		w.WriteHeader(http.StatusUnauthorized)
	case 403:
		w.WriteHeader(http.StatusForbidden)
	case 404:
		w.WriteHeader(http.StatusNotFound)
	case 429:
		w.WriteHeader(http.StatusTooManyRequests)
	case 500:
		w.WriteHeader(http.StatusInternalServerError)
	case 503:
		w.WriteHeader(http.StatusServiceUnavailable)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
