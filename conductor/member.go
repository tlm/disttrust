package conductor

import (
	"context"
	"time"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	logu "github.com/tlmiller/disttrust/log"
	"github.com/tlmiller/disttrust/provider"
)

type Member interface {
	DoneCh() <-chan error
	Play()
	Stop()
}

type DefaultMember struct {
	context  context.Context
	cancel   func()
	doneCh   chan error
	handle   LeaseHandler
	provider provider.Provider
	request  provider.Request
}

func (m *DefaultMember) DoneCh() <-chan error {
	return m.doneCh
}

func NewMember(provider provider.Provider, request provider.Request, handle LeaseHandler) Member {
	mem := &DefaultMember{
		doneCh:   make(chan error),
		handle:   handle,
		provider: provider,
		request:  request,
	}
	mem.context, mem.cancel = context.WithCancel(context.Background())
	return mem
}

func (m *DefaultMember) Play() {
	select {
	case <-m.context.Done():
		m.doneCh <- nil
		return
	default:
	}

	plog := log.WithFields(log.Fields{
		"common_name": m.request.CommonName,
	})

	plog.Info("acquiring lease")
	lease, err := m.provider.Issue(&m.request)
	if err != nil {
		m.doneCh <- errors.Wrap(err, "issuing certificate")
		return
	}

	plog = plog.WithFields(log.Fields{
		"lease_id":  lease.ID(),
		"lease_end": lease.Till().Format(time.RFC3339),
	})
	err = m.handle.Handle(logu.WithLogger(m.context, plog), lease)
	if err != nil {
		m.doneCh <- errors.Wrap(err, "handling new lease")
		return
	}

	for {
		plog = plog.WithFields(log.Fields{
			"lease_id":  lease.ID(),
			"lease_end": lease.Till().Format(time.RFC3339),
		})
		plog.Info("sleeping for next release renewal")
		r := NewRenewer(lease)
		go r.Renew()

		select {
		case <-r.RenewCh():
			plog.Info("renewing lease")
			lease, err = m.renewLease(lease)
			if err != nil {
				m.doneCh <- errors.Wrap(err, "renewing lease")
				return
			}
			plog.Info("acquired new lease")
			err = m.handle.Handle(logu.WithLogger(m.context, plog), lease)
			if err != nil {
				m.doneCh <- errors.Wrap(err, "handling renewed lease")
				return
			}

		case err = <-r.DoneCh():
			if err != nil {
				m.cancel()
				m.doneCh <- errors.Wrap(err, "waiting renewer")
				return
			}
		case <-m.context.Done():
			plog.Info("stopping member lease lifecycle")
			r.Stop()
			m.doneCh <- nil
			return
		}
	}
}

func (m *DefaultMember) renewLease(lease provider.Lease) (provider.Lease, error) {
	return nil, nil
}

func (m *DefaultMember) Stop() {
	m.cancel()
}
