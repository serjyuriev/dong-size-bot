package dongrepository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	c "github.com/serjyuriev/dong-size-bot/internal/pkg/config"
	m "github.com/serjyuriev/dong-size-bot/internal/pkg/models"
	_ "modernc.org/sqlite"
)

var (
	ErrNoSize = errors.New("no dong size")
)

type DongRepository interface {
	CreateDongSize(ctx context.Context, s *m.Dong) error
	ReadTodayDongSize(ctx context.Context, user int64) (*m.Dong, error)
	ReadUserAllTimeAverageSize(ctx context.Context, user int64) (float32, error)
	ReadUserCurrentYearAverageSize(ctx context.Context, user int64) (float32, error)
	ReadUserCurrentMonthAverageSize(ctx context.Context, user int64) (float32, error)
}

func CreateDongRepository(ctx context.Context) (DongRepository, error) {
	cfg, err := c.GetConfig()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", cfg.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &dongRepository{
		db: db,
	}, nil
}

type dongRepository struct {
	db *sql.DB
}

func (r *dongRepository) CreateDongSize(ctx context.Context, s *m.Dong) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		"insert into schlong_size(id, year, month, day, size) values ($1, $2, $3, $4, $5);",
		s.OwnerID,
		s.GeneratedAt.Year(),
		s.GeneratedAt.Month(),
		s.GeneratedAt.Day(),
		s.Size,
	); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *dongRepository) ReadTodayDongSize(ctx context.Context, user int64) (*m.Dong, error) {
	today := time.Now()
	row := r.db.QueryRowContext(
		ctx,
		"select id, size from schlong_size where year = $1 and month = $2 and day = $3 and id = $4;",
		today.Year(),
		today.Month(),
		today.Day(),
		user,
	)

	d := new(m.Dong)
	if err := row.Scan(&d.OwnerID, &d.Size); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoSize
		}
		return nil, err
	}

	return d, nil
}

func (r *dongRepository) ReadUserAllTimeAverageSize(ctx context.Context, user int64) (float32, error) {
	row := r.db.QueryRowContext(
		ctx,
		"select avg(size) from schlong_size where id = $1;",
		user,
	)

	var size float32
	if err := row.Scan(&size); err != nil {
		return 0, err
	}

	return size, nil
}

func (r *dongRepository) ReadUserCurrentYearAverageSize(ctx context.Context, user int64) (float32, error) {
	row := r.db.QueryRowContext(
		ctx,
		"select avg(size) from schlong_size where id = $1 and year = $2;",
		user,
		time.Now().Year(),
	)

	var size float32
	if err := row.Scan(&size); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoSize
		}
		return 0, err
	}

	return size, nil
}

func (r *dongRepository) ReadUserCurrentMonthAverageSize(ctx context.Context, user int64) (float32, error) {
	today := time.Now()
	row := r.db.QueryRowContext(
		ctx,
		"select avg(size) from schlong_size where id = $1 and year = $2 and month = $3;",
		user,
		today.Year(),
		today.Month(),
	)

	var size float32
	if err := row.Scan(&size); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoSize
		}
		return 0, err
	}

	return size, nil
}
