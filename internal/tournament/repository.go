//go:generate mockgen --source=repository.go -destination=repository_mock.go -package=tournament
package tournament

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const ErrCodeMySQLDuplicateEntry uint16 = 1062

var (
	ErrNotFound       = errors.New("not found")
	ErrDuplicateEntry = errors.New("already exists")
	ErrNotDeleted     = errors.New("not deleted")
	ErrNotUpdated     = errors.New("not updated")
)

type Repository interface {
	GetTournamen(ctx context.Context, id uint) (*Tournament, error)
	CreateTournament(ctx context.Context, tournament *Tournament) error
	UpdateTournament(ctx context.Context, tournament *Tournament) error
	DeleteTournament(ctx context.Context, tournament *Tournament) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (*RepositoryImpl) CreateTournament(ctx context.Context, tournament *Tournament) error {
	// TODO
	panic("unimplemented")
}

func (*RepositoryImpl) DeleteTournament(ctx context.Context, tournament *Tournament) error {
	// TODO
	panic("unimplemented")
}

func (*RepositoryImpl) GetTournamen(ctx context.Context, id uint) (*Tournament, error) {
	// TODO
	panic("unimplemented")
}
func (*RepositoryImpl) UpdateTournament(ctx context.Context, tournament *Tournament) error {
	// TODO
	panic("unimplemented")
}
