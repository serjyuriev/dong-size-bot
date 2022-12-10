package dongsizebot

import (
	"context"
	"time"

	cr "github.com/serjyuriev/dong-size-bot/internal/app/controller"
	dr "github.com/serjyuriev/dong-size-bot/internal/app/dong_repository"
	ds "github.com/serjyuriev/dong-size-bot/internal/app/dong_service"
	cf "github.com/serjyuriev/dong-size-bot/internal/pkg/config"
	tele "gopkg.in/telebot.v3"
)

type DongSizeBot interface {
	Start()
}

func NewDongSizeBot(ctx context.Context) (DongSizeBot, error) {
	cfg, err := cf.GetConfig()
	if err != nil {
		return nil, err
	}

	repo, err := dr.CreateDongRepository(ctx)
	if err != nil {
		return nil, err
	}

	svc, err := ds.CreateDongService(ctx, repo)
	if err != nil {
		return nil, err
	}

	ctrl, err := cr.CreateController(ctx, repo, svc)
	if err != nil {
		return nil, err
	}

	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &dongSizeBot{
		ctrl: ctrl,
		tb:   b,
	}, nil
}

type dongSizeBot struct {
	ctrl cr.Controller
	tb   *tele.Bot
}

func (b *dongSizeBot) Start() {
	b.tb.Handle("/cockannual", b.ctrl.AnnualDongSize)
	b.tb.Handle("/cocklifetime", b.ctrl.LifetimeDongSize)
	b.tb.Handle("/cockmonth", b.ctrl.MonthDongSize)
	b.tb.Handle("/cocksize", b.ctrl.TodayDongSize)

	b.tb.Start()
}
