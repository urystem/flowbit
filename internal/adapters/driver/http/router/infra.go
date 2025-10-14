package router

import (
	"net/http"
)

// type router struct {
// 	middleware inbound.MiddleWareInter
// 	handler    inbound.HandlerInter
// }

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

func NewRoute(hand HandlerInter) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", hand.Catalog)
	// mux.Handle("GET /", middle.CheckOrSetSession(http.HandlerFunc(hand.Catalog)))
	mux.HandleFunc("GET /postimage/{image}", hand.ServePostImage)
	mux.HandleFunc("GET /create-post-page", hand.CreatePostPage)
	mux.HandleFunc("POST /submit-post", hand.SubmitPost)
	mux.HandleFunc("GET /archive", hand.Archive)
	mux.HandleFunc("GET /archive-post/{PostID}", hand.ArchivePost)
	mux.HandleFunc("GET /comment/{image}", hand.ServeCommentImage)
	mux.HandleFunc("GET /post/{postID}", hand.ActivePost)
	mux.HandleFunc("POST /add-comment/{postID}", hand.AddComment)
	mux.HandleFunc("POST /reply/{commentID}", hand.Reply)

	return mux
	// return mux
}
