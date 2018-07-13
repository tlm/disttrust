package conductor

import (
	"context"
	"time"

	"github.com/carlescere/goback"

	"github.com/sirupsen/logrus"
)

type Conductor struct {
	cancel  func()
	context context.Context
	members []*Member
	watchCh chan error
}

const (
	MaxRetry int = 4
)

var (
	MaxBackoff = 30 * time.Second
	MinBackoff = 3 * time.Second
)

func NewConductor() *Conductor {
	return &Conductor{
		watchCh: make(chan error),
	}
}

func (c *Conductor) AddMember(mem Member) *Conductor {
	go func() {
		var err, retryErr error
		backoff := goback.SimpleBackoff{
			MaxAttempts: MaxRetry,
			Min:         MinBackoff,
			Max:         MaxBackoff,
		}
		for ; retryErr == nil; retryErr = goback.Wait(&backoff) {
			log := logrus.WithFields(logrus.Fields{
				"retries": backoff.Attempts,
			})
			log.Info("playing member")
			err = nil
			go mem.Play()

			select {
			case err = <-mem.DoneCh():
				if err == nil {
					log.Info("member play stopped")
					c.watchCh <- nil
					return
				}
				log.Errorf("member play stoped: %v", err)
			}
		}
		c.watchCh <- err
	}()
	return c
}

func (c *Conductor) Watch() error {
	for {
		select {
		case err := <-c.watchCh:
			return err
		}
	}
}
