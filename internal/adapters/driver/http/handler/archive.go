package handler

import (
	"log/slog"
	"net/http"

	"marketflow/internal/domain"
)

func (h *handler) Archive(w http.ResponseWriter, r *http.Request) {
	posts, err := h.use.ListOfArchivePosts(r.Context())
	if err != nil {
		slog.Error(err.Error())

		errData := &domain.ErrorPageData{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}

		h.renderError(w, errData)
		return
	}

	err = h.templates.ExecuteTemplate(w, "archive.html", posts)
	if err != nil {
		slog.Error(err.Error())
		errData := &domain.ErrorPageData{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		h.renderError(w, errData)
	}
}
