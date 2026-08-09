package infra

import (
	"MyMessenger/services/gateway/internal/config"
	"context"
	"sync"
	"time"

	"go.uber.org/fx"
)

type AutoCleaningRepo struct {
	mu sync.RWMutex
	m  map[any]time.Time

	tickerTiming time.Duration
}

func NewAutoCleaningRepo(lf fx.Lifecycle, conf *config.InfraConfig) *AutoCleaningRepo {
	repo := &AutoCleaningRepo{
		m:            make(map[any]time.Time),
		tickerTiming: conf.TickerCleanerTiming,
	}

	ctxC, ctxCancel := context.WithCancel(context.Background())
	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go repo.cleanMap(ctxC)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctxCancel()
			return nil
		},
	})

	return repo
}

func (r *AutoCleaningRepo) cleanMap(ctx context.Context) {
	ticker := time.NewTicker(r.tickerTiming)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var toDel []any

			r.mu.RLock()
			now := time.Now()
			for id, val := range r.m {
				if val.Before(now) {
					toDel = append(toDel, id)
				}
			}
			r.mu.RUnlock()

			if len(toDel) == 0 {
				continue
			}

			r.mu.Lock()
			now = time.Now()
			for _, id := range toDel {
				val, ok := r.m[id]
				if ok && val.Before(now) {
					delete(r.m, id)
				}
			}
			r.mu.Unlock()
		}
	}
}

func (r *AutoCleaningRepo) Put(key any, expAt time.Time) {
	r.mu.Lock()
	r.m[key] = expAt
	r.mu.Unlock()
}

func (r *AutoCleaningRepo) Get(key any) (time.Time, bool) {
	r.mu.RLock()
	val, isExists := r.m[key]
	r.mu.RUnlock()

	return val, isExists
}

func (r *AutoCleaningRepo) IsExists(key any) bool {
	r.mu.RLock()
	_, isExists := r.m[key]
	r.mu.RUnlock()

	return isExists
}
