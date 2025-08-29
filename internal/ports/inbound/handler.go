package inbound

import "net/http"

type HandlerInter interface {
	Catalog(w http.ResponseWriter, r *http.Request)
	ServePostImage(w http.ResponseWriter, r *http.Request)
	CreatePostPage(w http.ResponseWriter, r *http.Request)
	SubmitPost(w http.ResponseWriter, r *http.Request)
	Archive(w http.ResponseWriter, r *http.Request)
	ArchivePost(w http.ResponseWriter, r *http.Request)
	ServeCommentImage(w http.ResponseWriter, r *http.Request)
	ActivePost(w http.ResponseWriter, r *http.Request)
	AddComment(w http.ResponseWriter, r *http.Request)
	Reply(w http.ResponseWriter, r *http.Request)
	// CreateComment(w http.ResponseWriter, r *http.Request)
	// ErrorHandler
}

// it is for middleware
// but we cannot give to middleware6 because handler is depency middleware
// type ErrorHandler interface {
// 	RenderError(w http.ResponseWriter, errPage *domain.ErrorPageData)
// }
