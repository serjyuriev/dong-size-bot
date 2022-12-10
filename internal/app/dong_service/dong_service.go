package dongservice

import (
	"context"
	"errors"
	"math/rand"
	"time"

	r "github.com/serjyuriev/dong-size-bot/internal/app/dong_repository"
	c "github.com/serjyuriev/dong-size-bot/internal/pkg/config"
	m "github.com/serjyuriev/dong-size-bot/internal/pkg/models"
)

type DongService interface {
	GetTodayDongSize(ctx context.Context, user int64) (*m.Dong, error)
}

func CreateDongService(ctx context.Context, dr r.DongRepository) (DongService, error) {
	if dr == nil {
		return nil, errors.New("dong repository is nil")
	}

	cfg, err := c.GetConfig()
	if err != nil {
		return nil, err
	}

	return &dongService{
		cfg: cfg,
		dr:  dr,
	}, nil
}

type dongService struct {
	cfg c.Config
	dr  r.DongRepository
}

func (s *dongService) GetTodayDongSize(ctx context.Context, user int64) (*m.Dong, error) {
	d, err := s.dr.ReadTodayDongSize(ctx, user)
	if err != nil {
		if errors.Is(err, r.ErrNoSize) {
			d = new(m.Dong)
			d.OwnerID = user
			d.GeneratedAt = time.Now()
			d.Size = s.GenerateDong(ctx)

			if err = s.dr.CreateDongSize(ctx, d); err != nil {
				return nil, err
			}

			return d, nil
		}
		return nil, err
	}

	return d, nil
}

func (s *dongService) GenerateDong(ctx context.Context) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(s.cfg.MaxDongSize) + 1
}
