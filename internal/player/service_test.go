package player

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_GetPlayer_ShouldSucceed(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	var id uint = 1
	p := &Player{
		ID:     id,
		Name:   "John",
		Rating: 1000,
	}

	repo.EXPECT().GetPlayer(context.Background(), id).Return(p, nil)

	player, err := service.GetPlayer(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, p, player)
}

func Test_GetPlayer_ShouldFail_NotFound(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	var id uint = 1

	repo.EXPECT().GetPlayer(context.Background(), id).Return(nil, ErrNotFound)

	player, err := service.GetPlayer(context.Background(), id)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, player)
}

func Test_GetPlayer_ShouldFail(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	var (
		id      uint = 1
		errTest      = errors.New("test error")
	)

	repo.EXPECT().GetPlayer(context.Background(), id).Return(nil, errTest)

	player, err := service.GetPlayer(context.Background(), id)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errTest)
	assert.Nil(t, player)
}

func Test_GetPlayers_ShouldSucceed(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	var (
		id1 uint = 1
		id2 uint = 2
		p1       = &Player{
			ID:     id1,
			Name:   "John",
			Rating: 1000,
		}
		p2 = &Player{
			ID:     id2,
			Name:   "Jane",
			Rating: 1000,
		}
		ids         = []uint{id1, id2}
		testPlayers = []*Player{p1, p2}
	)

	repo.EXPECT().GetPlayers(context.Background(), ids).Return(testPlayers, nil)

	players, err := service.GetPlayers(context.Background(), ids)

	assert.NoError(t, err)
	assert.Equal(t, testPlayers, players)
}

func Test_GetPlayers_ShouldFail(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	ids := []uint{1, 2}
	errTest := errors.New("test error")

	repo.EXPECT().GetPlayers(context.Background(), ids).Return(nil, errTest)

	players, err := service.GetPlayers(context.Background(), ids)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errTest)
	assert.Nil(t, players)
}

func Test_GetPlayers_ShouldFail_NotFound(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	ids := []uint{1, 2}

	repo.EXPECT().GetPlayers(context.Background(), ids).Return(nil, ErrNotFound)

	players, err := service.GetPlayers(context.Background(), ids)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, players)
}

func Test_GetTopPlayersByRating_ShouldSucced(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	var (
		id1 uint = 1
		id2 uint = 2
		p1       = &Player{
			ID:     id1,
			Name:   "John",
			Rating: 1000,
		}
		p2 = &Player{
			ID:     id2,
			Name:   "Jane",
			Rating: 1000,
		}
		testPlayers = []*Player{p1, p2}
	)

	repo.EXPECT().GetTopPlayersByRating(context.Background(), 2).Return(testPlayers, nil)

	players, err := service.GetTopPlayersByRating(context.Background(), 2)

	assert.NoError(t, err)
	assert.Equal(t, testPlayers, players)
}

func Test_GetTopPlayersByRating_ShouldFail(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	errTest := errors.New("test error")

	repo.EXPECT().GetTopPlayersByRating(context.Background(), 2).Return(nil, errTest)

	players, err := service.GetTopPlayersByRating(context.Background(), 2)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errTest)
	assert.Nil(t, players)
}

func Test_CreatePlayer_ShouldSucceed(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	name := "John"
	player := &Player{
		ID:     0,
		Name:   name,
		Rating: 1000,
	}

	repo.EXPECT().CreatePlayer(context.Background(), player).Return(nil)

	err := service.CreatePlayer(context.Background(), name)

	assert.NoError(t, err)
}

func Test_CreatePlayer_ShouldFail(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	name := "John"
	player := &Player{
		ID:     0,
		Name:   name,
		Rating: 1000,
	}
	errTest := errors.New("test error")

	repo.EXPECT().CreatePlayer(context.Background(), player).Return(errTest)

	err := service.CreatePlayer(context.Background(), name)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errTest)
}

func Test_DeletePlayer_ShouldSucceed(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	var id uint = 1

	repo.EXPECT().DeletePlayer(context.Background(), id).Return(nil)

	err := service.DeletePlayer(context.Background(), id)

	assert.NoError(t, err)
}

func Test_DeletePlayer_ShouldFail(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	var id uint = 1
	errTest := errors.New("test error")

	repo.EXPECT().DeletePlayer(context.Background(), id).Return(errTest)

	err := service.DeletePlayer(context.Background(), id)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errTest)
}

func Test_UpdatePlayers_ShouldSucceed(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	var (
		p1 = &Player{
			ID:     1,
			Name:   "John",
			Rating: 1000,
		}
		p2 = &Player{
			ID:     2,
			Name:   "Jane",
			Rating: 1000,
		}
		players    = []*Player{p1, p2}
		newRatings = []int{1200, 800}
	)

	repo.EXPECT().UpdatePlayers(context.Background(), players).Return(nil)

	err := service.UpdatePlayers(context.Background(), players, newRatings)

	assert.NoError(t, err)
}

func Test_UpdatePlayers_Shouldfail(t *testing.T) {
	t.Parallel()

	repo := NewMockRepository(gomock.NewController(t))
	service := NewService(repo)

	var (
		p1 = &Player{
			ID:     1,
			Name:   "John",
			Rating: 1000,
		}
		p2 = &Player{
			ID:     2,
			Name:   "Jane",
			Rating: 1000,
		}
		players    = []*Player{p1, p2}
		newRatings = []int{1200, 800}
		errTest    = errors.New("test error")
	)

	repo.EXPECT().UpdatePlayers(context.Background(), players).Return(errTest)

	err := service.UpdatePlayers(context.Background(), players, newRatings)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errTest)
}
