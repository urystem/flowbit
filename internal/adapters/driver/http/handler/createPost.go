package handler

// import (
// 	"log/slog"
// 	"net/http"

// 	"marketflow/internal/domain"
// )

// func (h *handler) CreatePostPage(w http.ResponseWriter, r *http.Request) {
// 	h.templates.ExecuteTemplate(w, "create-post.html", nil)
// }

// func (h *handler) SubmitPost(w http.ResponseWriter, r *http.Request) {
// 	// Парсим multipart форму
// 	// err = r.ParseMultipartForm(10 << 20) // 10 MB
// 	// if err != nil {
// 	// 	http.Error(w, "Unable to parse form", http.StatusBadRequest)
// 	// 	return
// 	// }
// 	ctx := r.Context()
// 	form := &domain.Form{
// 		// Name:    r.FormValue("name"),
// 		// Subject: r.FormValue("subject"),
// 		// Content: r.FormValue("comment"),
// 	}
// 	form.Name = r.FormValue("name")
// 	form.Subject = r.FormValue("subject")
// 	form.Content = r.FormValue("content")

// 	file, header, err := r.FormFile("file")
// 	if err == nil {
// 		defer file.Close()
// 		form.File = new(domain.InPutObject)
// 		form.File.Reader = file
// 		form.File.ConType = header.Header.Get("Content-Type")
// 		form.File.Size = header.Size
// 	} else if err != http.ErrMissingFile {
// 		http.Error(w, "File error", http.StatusInternalServerError)
// 		slog.Error(err.Error())
// 		errData := &domain.ErrorPageData{
// 			Code:    http.StatusInternalServerError,
// 			Message: err.Error(),
// 		}

// 		h.renderError(w, errData)
// 		return
// 	}

// 	err = h.use.CreatePost(ctx, form)
// 	if err != nil {
// 		errData := &domain.ErrorPageData{
// 			Code:    http.StatusInternalServerError,
// 			Message: err.Error(),
// 		}
// 		h.renderError(w, errData)
// 		return
// 	}
// 	h.Catalog(w, r)
// }
