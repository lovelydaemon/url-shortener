package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/pkg/logger"
	"github.com/lovelydaemon/url-shortener/internal/storage"
)

type Queue struct {
	c chan entity.UserWithURLs
	s storage.Storage
	l logger.Interface
}

// New creates a new queue
func New(s storage.Storage, l logger.Interface) *Queue {
	return &Queue{
		c: make(chan entity.UserWithURLs, 1),
		s: s,
		l: l,
	}
}

func (q *Queue) Push(value entity.UserWithURLs) {
	q.c <- value
}

func (q *Queue) FlushUserURLs(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)

	var shortURLs []entity.UserWithURLs

	for {
		select {
		case str := <-q.c:
			shortURLs = append(shortURLs, str)
		case <-ctx.Done():
			ticker.Stop()
			close(q.c)
			return
		case <-ticker.C:
			if len(shortURLs) == 0 {
				continue
			}

			if err := q.s.DeleteURLs(ctx, shortURLs...); err != nil {
				q.l.Error(fmt.Errorf("Queue - FlushUserURLs - q.s.DeleteURLs: %w", err))
				continue
			}

			shortURLs = nil
		}
	}
}
