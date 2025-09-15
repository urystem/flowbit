package outbound

// import (
// 	"context"

// 	"marketflow/internal/domain"

// 	"github.com/google/uuid"
// )

// type PostGres interface {
// 	Pgx       // for usecase
// 	CloseDB() // for bootstrap
// }

// type Pgx interface {
// 	pgxPost
// 	pgxUser
// 	pgxComment
// }

// type pgxPost interface {
// 	SelectActivePosts(context.Context) ([]domain.PostNonContent, error)
// 	InsertPost(ctx context.Context, post *domain.InsertPost) (uint64, error)
// 	DeletePost(ctx context.Context, id uint64) error
// 	Archiver(ctx context.Context) error
// 	SelectArchivePosts(context.Context) ([]domain.PostNonContent, error)
// 	GetPost(ctx context.Context, id uint64) (*domain.PostX, error)
// }

// type pgxUser interface {
// 	InsertUser(context.Context, *domain.Session) error
// 	DeleteUser(ctx context.Context, sessionID uuid.UUID) error
// }

// type pgxComment interface {
// 	GetComments(ctx context.Context, postID uint64) ([]domain.Comment, error)
// 	InsertComment(ctx context.Context, comment *domain.InsertComment) (uint64, error)
// 	DeleteComment(ctx context.Context, commentID uint64) error
// 	InsertReply(ctx context.Context, reply *domain.InsertReply) (uint64, error)
// }
