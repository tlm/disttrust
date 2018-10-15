package healthz

import (
	"sync"
)

type Healthz struct {
	checks []Checker
	lock   sync.Mutex
}

func (h *Healthz) Checks() []Checker {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.checks
}

func New() *Healthz {
	return &Healthz{
		checks: []Checker{},
		lock:   sync.Mutex{},
	}
}

func (h *Healthz) SetChecks(checks ...Checker) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.checks = checks
}
