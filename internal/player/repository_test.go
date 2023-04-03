package player

import (
	"context"
	"foosball/internal/mysql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	goMySql "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	queryGetPlayer             = regexp.QuoteMeta("SELECT * FROM `players` WHERE `players`.`id` = ? ORDER BY `players`.`id` LIMIT 1")
	queryGetPlayers            = regexp.QuoteMeta("SELECT * FROM `players` WHERE `players`.`id` IN (?,?)")
	queryCreatePlayer          = regexp.QuoteMeta("INSERT INTO `players` (`name`,`rating`,`created_at`) VALUES (?,?,?)")
	queryDeletePlayer          = regexp.QuoteMeta("DELETE FROM `players` WHERE `players`.`id` = ?")
	queryGetTopPlayersByRating = regexp.QuoteMeta("SELECT * FROM `players` ORDER BY rating DESC LIMIT 1")
)

func Test_Repository_GetPlayer_ShouldSucceed(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)

	sql.
		ExpectQuery(queryGetPlayer).
		WithArgs(1).
		WillReturnRows(
			sql.NewRows([]string{"id", "name", "rating"}).
				AddRow(1, "Player 1", 1000),
		)

	player, err := repo.GetPlayer(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, &Player{
		ID:     1,
		Name:   "Player 1",
		Rating: 1000,
	}, player)
}

func Test_Repository_GetPlayer_ShouldFail_RecordNotFound(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)

	sql.
		ExpectQuery(queryGetPlayer).
		WithArgs(1).
		WillReturnError(gorm.ErrRecordNotFound)

	player, err := repo.GetPlayer(context.Background(), 1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, player)
}

func Test_Repository_GetPlayer_ShouldFail_Error(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)
	errTest := errors.New("test error")

	sql.
		ExpectQuery(queryGetPlayer).
		WithArgs(1).
		WillReturnError(errTest)

	player, err := repo.GetPlayer(context.Background(), 1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errTest)
	assert.Nil(t, player)
}

func Test_Repository_GetPlayers_ShouldSucceed(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)

	sql.
		ExpectQuery(queryGetPlayers).
		WithArgs(1, 2).
		WillReturnRows(
			sql.NewRows([]string{"id", "name", "rating"}).
				AddRow(1, "Player 1", 1000).
				AddRow(2, "Player 2", 1000),
		)

	players, err := repo.GetPlayers(context.Background(), []uint{1, 2})

	assert.NoError(t, err)
	assert.Equal(t, []*Player{
		{
			ID:     1,
			Name:   "Player 1",
			Rating: 1000,
		},
		{
			ID:     2,
			Name:   "Player 2",
			Rating: 1000,
		},
	}, players)
}

func Test_Repository_GetPlayers_ShouldFail_RecordNotFound(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)

	sql.
		ExpectQuery(queryGetPlayers).
		WithArgs(1, 2).
		WillReturnError(gorm.ErrRecordNotFound)

	players, err := repo.GetPlayers(context.Background(), []uint{1, 2})

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, players)
}

func Test_Repository_GetPlayers_ShouldFail_Error(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)
	errTest := errors.New("test error")

	sql.
		ExpectQuery(queryGetPlayers).
		WithArgs(1, 2).
		WillReturnError(errTest)

	players, err := repo.GetPlayers(context.Background(), []uint{1, 2})

	assert.Error(t, err)
	assert.ErrorIs(t, err, errTest)
	assert.Nil(t, players)
}

func Test_Repository_CreatePlayer_ShouldSucceed(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)

	sql.ExpectBegin()
	sql.ExpectExec(queryCreatePlayer).
		WithArgs("Player 1", 1000, mysql.AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sql.ExpectCommit()

	player := &Player{
		Name:      "Player 1",
		Rating:    1000,
		CreatedAt: time.Now(),
	}

	err := repo.CreatePlayer(context.Background(), player)

	assert.NoError(t, err)
}

func Test_Repository_CreatePlayer_ShouldFail_Duplicate(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)

	sql.ExpectBegin()
	sql.ExpectExec(queryCreatePlayer).
		WithArgs("Player 1", 1000, mysql.AnyTime{}).
		WillReturnError(&goMySql.MySQLError{Number: 1062})
	sql.ExpectRollback()

	player := &Player{
		Name:      "Player 1",
		Rating:    1000,
		CreatedAt: time.Now(),
	}

	err := repo.CreatePlayer(context.Background(), player)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrDuplicateEntry)
}

func Test_Repository_CreatePlayer_ShouldFail_Error(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)
	errTest := errors.New("test error")

	sql.ExpectBegin()
	sql.ExpectExec(queryCreatePlayer).
		WithArgs("Player 1", 1000, mysql.AnyTime{}).
		WillReturnError(errTest)
	sql.ExpectRollback()

	player := &Player{
		Name:      "Player 1",
		Rating:    1000,
		CreatedAt: time.Now(),
	}

	err := repo.CreatePlayer(context.Background(), player)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errTest)
}

func Test_Repository_DeletePlayer_ShouldSucceed(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)

	sql.ExpectBegin()
	sql.ExpectExec(queryDeletePlayer).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sql.ExpectCommit()

	err := repo.DeletePlayer(context.Background(), 1)

	assert.NoError(t, err)
}

func Test_Repository_DeletePlayer_ShouldFail_NoneDeleted(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)

	sql.ExpectBegin()
	sql.ExpectExec(queryDeletePlayer).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 0))
	sql.ExpectCommit()

	err := repo.DeletePlayer(context.Background(), 1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNoneDeleted)
}

func Test_Repository_DeletePlayer_ShouldFail_Error(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)
	errTest := errors.New("test error")

	sql.ExpectBegin()
	sql.ExpectExec(queryDeletePlayer).
		WithArgs(1).
		WillReturnError(errTest)
	sql.ExpectRollback()

	err := repo.DeletePlayer(context.Background(), 1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errTest)
}

func Test_Repository_UpdatePlayer_ShouldSucceed(t *testing.T) {
	// TODO
	assert.Nil(t, 1)
}

func Test_Repository_UpdatePlayer_ShouldFail_NotUpdated(t *testing.T) {
	// TODO
	assert.Nil(t, 1)
}

func Test_Repository_UpdatePlayer_ShouldFail_Error(t *testing.T) {
	// TODO
	assert.Nil(t, 1)
}

func Test_Repository_GetTopPlayersByRating_ShouldSucceed(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)

	sql.ExpectQuery(queryGetTopPlayersByRating).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1).
			AddRow(2))

	players, err := repo.GetTopPlayersByRating(context.Background(), 2)

	assert.NoError(t, err)
	assert.Len(t, players, 2)
}

func Test_Repository_GetTopPlayersByRating_ShouldFail(t *testing.T) {
	t.Parallel()

	db, sql := mysql.NewMockClient(t)
	repo := NewRepository(db)
	errTest := errors.New("test error")

	sql.ExpectQuery(queryGetTopPlayersByRating).
		WillReturnError(errTest)
	sql.ExpectRollback()

	players, err := repo.GetTopPlayersByRating(context.Background(), 1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errTest)
	assert.Nil(t, players)
}
