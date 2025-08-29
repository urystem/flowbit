package inbound

import (
	"context"

	"1337b04rd/internal/domain"
)

type UseCase interface {
	Service // for handler
	Ticker  // for bootstrap
}

type Service interface {
	postUse
	commentUsecase
	userUseCase
}

type postUse interface {
	GetPostImage(ctx context.Context, objName string) (*domain.OutputObject, error)
	CreatePost(ctx context.Context, form *domain.Form) error
	activePostUse
	archivePostUse
}

type activePostUse interface {
	ListOfActivePosts(context.Context) ([]domain.PostNonContent, error)
	GetActivePost(context.Context, uint64) (*domain.ActivePost, error)
}

type archivePostUse interface {
	ListOfArchivePosts(context.Context) ([]domain.PostNonContent, error)
	GetArchivePost(context.Context, uint64) (*domain.ArchivePost, error)
}

type commentUsecase interface {
	CreateComment(ctx context.Context, form *domain.CommentForm) error
	GetCommentImage(ctx context.Context, objName string) (*domain.OutputObject, error)
	Reply(context.Context, *domain.ReplyForm) error
}

type userUseCase interface {
	AddUserToDB(ctx context.Context, ses *domain.Session) error
}

type Ticker interface {
	Archiver(ctx context.Context) error
}
