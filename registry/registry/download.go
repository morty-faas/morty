package registry

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) GetImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id") + "-" + chi.URLParam(r, "version")

	_, link, err := s.storage.GetDownloadLink(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, link, http.StatusFound)
}
