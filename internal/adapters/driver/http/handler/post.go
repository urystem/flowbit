package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"marketflow/internal/domain"
)

func (h *handler) ActivePost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(r.PathValue("postID"), 10, 64)
	if err != nil {
		slog.Error(err.Error())

		errData := &domain.ErrorPageData{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}

		h.renderError(w, errData)
		return
	}

	post, err := h.use.GetActivePost(r.Context(), postID)
	if err != nil {
		slog.Error(err.Error())

		errData := &domain.ErrorPageData{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}

		h.renderError(w, errData)
		return
	}

	err = h.templates.ExecuteTemplate(w, "post.html", post)
	if err != nil {
		slog.Error(err.Error())
		errData := &domain.ErrorPageData{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		h.renderError(w, errData)
	}
}
