package router

import (
	"net/http"

	"1337b04rd/internal/ports/inbound"
)

// type router struct {
// 	middleware inbound.MiddleWareInter
// 	handler    inbound.HandlerInter
// }

func NewRoute(middle inbound.MiddleWareInter, hand inbound.HandlerInter) http.Handler {
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

	return middle.CheckOrSetSession(mux)
	// return mux
}
