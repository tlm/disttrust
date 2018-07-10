package conductor

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/tlmiller/disttrust/provider"
)

// Type used for handling then when of a provider lease. Does not do any renewal
// of lease but instead act as the signaller for when a lease should be renewed.
// Inspired by https://github.com/hashicorp/vault/blob/master/api/renewer.go#L47
type Renewer struct {
	lease   provider.Lease
	lock    sync.Mutex
	doneCh  chan error
	renewCh chan struct{}
	stopCh  chan bool
	stopped bool
}

var (
	ErrLeaseLapsed = errors.New("lease has already lapsed")
	// The threshold of time we don't even bother about waiting for. Anything
	// below this mark will trigger straight away.
	RenewThreshold = 10 * time.Second
)

// Cqlculates a grace period after 80% of the lease duration so that we create
// a random window within 80-90& of the lease period. The returned value is a
// duration from 80% that should be waited on. If the lease duration is 0 then
// the grace period is automatically zero
func calculateGrace(leaseDuration time.Duration) time.Duration {
	if leaseDuration == 0 || leaseDuration < 10*time.Nanosecond {
		return 0
	}

	rand := rand.New(rand.NewSource(time.Now().Unix()))
	leaseNanos := float64(leaseDuration.Nanoseconds())
	// Work out 10% of the lease duration and that is our maximum jitter
	jitterMax := leaseNanos * 0.1
	return time.Duration(uint64(rand.Int63()) % uint64(jitterMax))
}

func (r *Renewer) DoneCh() <-chan error {
	return r.doneCh
}

func NewRenewer(lease provider.Lease) *Renewer {
	return &Renewer{
		lease:   lease,
		doneCh:  make(chan error, 1),
		renewCh: make(chan struct{}),
		stopCh:  make(chan bool),
		stopped: false,
	}
}

func (r *Renewer) Renew() {
	left := time.Until(r.lease.Till())
	if left <= time.Duration(0) {
		r.doneCh <- ErrLeaseLapsed
		return
	}
	if left <= RenewThreshold {
		r.renewCh <- struct{}{}
		return
	}

	// find out where 80% of the time left is
	minRenew := time.Duration(float64(left) * 0.8)
	timer := time.NewTimer(minRenew + calculateGrace(left))

	select {
	case <-timer.C:
		r.renewCh <- struct{}{}
		r.doneCh <- nil
		return
	case <-r.stopCh:
		r.doneCh <- nil
		return
	}
}

func (r *Renewer) RenewCh() <-chan struct{} {
	return r.renewCh
}

func (r *Renewer) Stop() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if !r.stopped {
		close(r.stopCh)
		r.stopped = true
	}
}
