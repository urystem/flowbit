package handler

// import (
// 	"log/slog"
// 	"net/http"
// 	"strconv"

// 	"marketflow/internal/domain"
// )

// func (h *handler) ArchivePost(w http.ResponseWriter, r *http.Request) {
// 	postID, err := strconv.ParseUint(r.PathValue("PostID"), 10, 64)
// 	if err != nil {
// 		slog.Error(err.Error())

// 		errData := &domain.ErrorPageData{
// 			Code:    http.StatusBadRequest,
// 			Message: err.Error(),
// 		}

// 		h.renderError(w, errData)
// 		return
// 	}

// 	post, err := h.use.GetArchivePost(r.Context(), postID)
// 	if err != nil {
// 		slog.Error(err.Error())

// 		errData := &domain.ErrorPageData{
// 			Code:    http.StatusInternalServerError,
// 			Message: err.Error(),
// 		}

// 		h.renderError(w, errData)
// 		return
// 	}

// 	err = h.templates.ExecuteTemplate(w, "archive-post.html", post)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		errData := &domain.ErrorPageData{
// 			Code:    http.StatusInternalServerError,
// 			Message: err.Error(),
// 		}
// 		h.renderError(w, errData)
// 	}
// }

// func (h *handler) ServeCommentImage(w http.ResponseWriter, r *http.Request) {
// 	// Получаем имя файла из URL
// 	imageName := r.PathValue("image")
// 	if imageName == "" {
// 		http.Error(w, "missing image name", http.StatusBadRequest)
// 		return
// 	}

// 	// Получаем объект из MinIO
// 	obj, err := h.use.GetCommentImage(r.Context(), imageName)
// 	if err != nil {
// 		slog.Error("", "get object:", err)
// 		http.Error(w, "file not found", http.StatusNotFound)
// 		return
// 	}
// 	defer obj.Close()

// 	w.Header().Set("Content-Type", obj.ConType)

// 	http.ServeContent(w, r, "", obj.Modified, obj)
// 	obj.Close()
// }
