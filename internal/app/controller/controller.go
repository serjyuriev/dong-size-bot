package controller

import (
	"context"
	"errors"
	"fmt"

	r "github.com/serjyuriev/dong-size-bot/internal/app/dong_repository"
	s "github.com/serjyuriev/dong-size-bot/internal/app/dong_service"
	tele "gopkg.in/telebot.v3"
)

type Controller interface {
	AnnualDongSize(c tele.Context) error
	LifetimeDongSize(c tele.Context) error
	MonthDongSize(c tele.Context) error
	TodayDongSize(c tele.Context) error
}

func CreateController(ctx context.Context, repo r.DongRepository, svc s.DongService) (Controller, error) {
	if repo == nil {
		return nil, errors.New("repository is nil")
	}

	if svc == nil {
		return nil, errors.New("service is nil")
	}

	return &controller{
		repo: repo,
		svc:  svc,
	}, nil
}

type controller struct {
	repo r.DongRepository
	svc  s.DongService
}

func (ctrl *controller) AnnualDongSize(c tele.Context) error {
	average, err := ctrl.repo.ReadUserCurrentYearAverageSize(context.Background(), c.Sender().ID)
	if err != nil {
		return c.Send("Случилась неприятность, зовите Сержа")
	}

	n := ctrl.getDisplayedName(c.Sender())
	e := ctrl.getEmoji(int(average))

	return c.Send(
		fmt.Sprintf(
			"%s, средний размер твоего члена в этом году - %.1f см! %s",
			n,
			average,
			e,
		),
	)
}

func (ctrl *controller) LifetimeDongSize(c tele.Context) error {
	average, err := ctrl.repo.ReadUserAllTimeAverageSize(context.Background(), c.Sender().ID)
	if err != nil {
		return c.Send("Случилась неприятность, зовите Сержа")
	}

	n := ctrl.getDisplayedName(c.Sender())
	e := ctrl.getEmoji(int(average))

	return c.Send(
		fmt.Sprintf(
			"%s, средний размер твоего члена за всю жизнь - %.1f см! %s",
			n,
			average,
			e,
		),
	)
}

func (ctrl *controller) MonthDongSize(c tele.Context) error {
	average, err := ctrl.repo.ReadUserCurrentMonthAverageSize(context.Background(), c.Sender().ID)
	if err != nil {
		return c.Send("Случилась неприятность, зовите Сержа")
	}

	n := ctrl.getDisplayedName(c.Sender())
	e := ctrl.getEmoji(int(average))

	return c.Send(
		fmt.Sprintf(
			"%s, средний размер твоего члена за текущий месяц - %.1f см! %s",
			n,
			average,
			e,
		),
	)
}

func (ctrl *controller) TodayDongSize(c tele.Context) error {
	ds, err := ctrl.svc.GetTodayDongSize(context.Background(), c.Sender().ID)
	if err != nil {
		return c.Send("Случилась неприятность, зовите Сержа")
	}

	n := ctrl.getDisplayedName(c.Sender())
	e := ctrl.getEmoji(ds.Size)

	return c.Send(
		fmt.Sprintf(
			"%s, размер твоего члена сегодня - %d см! %s",
			n,
			ds.Size,
			e,
		),
	)
}

func (ctrl *controller) getDisplayedName(u *tele.User) string {
	if u.Username == "" {
		return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	} else {
		return u.Username
	}
}

func (ctrl *controller) getEmoji(size int) string {
	if size == 35 {
		return "🥵"
	}

	if size > 29 {
		return "😱"
	}

	if size > 24 {
		return "😎"
	}

	if size > 19 {
		return "😳"
	}

	if size > 14 {
		return "🙂"
	}

	if size > 9 {
		return "😅"
	}

	if size > 4 {
		return "😬"
	}

	return "🥶"
}
