package conductor

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	"github.com/tlmiller/disttrust/provider"
)

type Member struct {
	doneCh   chan error
	lock     sync.Mutex
	provider provider.Provider
	request  provider.Request
	stopCh   chan struct{}
	stopped  bool
}

func (m *Member) DoneCh() <-chan error {
	return m.doneCh
}

func (m *Member) handleLease(lease provider.Lease) error {
	if !lease.HasResponse() {
		return nil
	}

	res, err := lease.Response()
	if err != nil {
		return errors.Wrap(err, "handling new lease response")
	}
	fmt.Println(res)

	return errors.New("test")
}

func NewMember(provider provider.Provider, request provider.Request) *Member {
	return &Member{
		doneCh:   make(chan error),
		provider: provider,
		request:  request,
		stopCh:   make(chan struct{}),
		stopped:  false,
	}
}

func (m *Member) Play() {
	select {
	case <-m.stopCh:
		m.doneCh <- nil
		return
	default:
	}

	log := log.WithFields(log.Fields{
		"common_name": m.request.CommonName,
	})

	log.Info("acquiring lease")
	lease, err := m.provider.Issue(&m.request)
	if err != nil {
		log.Error(err)
		m.doneCh <- errors.Wrap(err, "issuing certificate")
		return
	}

	err = m.handleLease(lease)
	if err != nil {
		m.doneCh <- errors.Wrap(err, "handling new lease")
		return
	}

	for {
		r := NewRenewer(lease)
		go r.Renew()

		select {
		case <-r.RenewCh():
			lease, err = m.renewLease(lease)
			if err != nil {
				m.doneCh <- errors.Wrap(err, "renewing lease")
				return
			}
			err = m.handleLease(lease)
			if err != nil {
				m.doneCh <- errors.Wrap(err, "handling renewed lease")
				return
			}

		case err = <-r.DoneCh():
			if err != nil {
				m.doneCh <- errors.Wrap(err, "waiting renewer")
			}
		case <-r.stopCh:
			m.doneCh <- nil
			return
		}
	}
}

func (m *Member) renewLease(lease provider.Lease) (provider.Lease, error) {
	return nil, nil
}

func (m *Member) Stop() {
	m.lock.Lock()
	defer m.lock.Unlock()
	if !m.stopped {
		close(m.stopCh)
		m.stopped = true
	}
}
