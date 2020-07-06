package main

import (
	"time"
)

// DailyCron repeats given func at specified time of day
type dailyCron struct {
	t          *time.Timer
	HourToTick int
	MinToTick  int
	SecToTick  int
	LastTick   time.Time
	loc        *time.Location
}

func newDailyCron(h, m, s int, tz string) (dailyCron, error) {
	tl, err := time.LoadLocation(tz)
	if err != nil {
		return dailyCron{}, err
	}
	dc := dailyCron{nil, h, m, s, time.Time{}, tl}
	dc.t = time.NewTimer(dc.getNextTickDuration())
	return dc, nil
}

func (d *dailyCron) updateTimer() {
	d.t.Reset(d.getNextTickDuration())
}

func (d *dailyCron) getNextTickDuration() time.Duration {
	now := time.Now().In(d.loc)
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), d.HourToTick, d.MinToTick, d.SecToTick, 0, d.loc)
	if now.After(nextTick) || d.LastTick.Day() == now.Day() {
		nextTick = nextTick.AddDate(0, 0, 1)
	}
	return nextTick.Sub(now)
}
