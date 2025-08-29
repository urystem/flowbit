package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"1337b04rd/internal/domain"
)

func (h *handler) Reply(w http.ResponseWriter, r *http.Request) {
	commentID, err := strconv.ParseUint(r.PathValue("commentID"), 10, 64)
	if err != nil {
		slog.Error(err.Error())

		errData := &domain.ErrorPageData{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}

		h.renderError(w, errData)
		return
	}

	ctx := r.Context()
	sess, x := h.middleware.FromContext(ctx)
	if !x {
		http.Error(w, "error middleware", http.StatusUnauthorized)
		return
	}
	form := &domain.ReplyForm{}
	form.User = sess.Uuid
	form.ReplyToID = commentID
	form.Content = r.FormValue("comment")
	file, header, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
		form.File = new(domain.InPutObject)
		form.File.Reader = file
		form.File.ConType = header.Header.Get("Content-Type")
		form.File.Size = header.Size
	} else if err != http.ErrMissingFile {
		slog.Error(err.Error())
		errData := &domain.ErrorPageData{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		h.renderError(w, errData)
		return
	}

	if !sess.Saved {
		err := h.use.AddUserToDB(ctx, sess)
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}

	err = h.use.Reply(ctx, form)
	if err != nil {
		slog.Error(err.Error())
		errData := &domain.ErrorPageData{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		h.renderError(w, errData)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
