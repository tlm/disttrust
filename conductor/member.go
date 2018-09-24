package conductor

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	logu "github.com/tlmiller/disttrust/log"
	"github.com/tlmiller/disttrust/provider"
)

type Member interface {
	DoneCh() <-chan error
	Name() string
	Play()
	Stop()
}

// DefaultMember is responsible for driving the main cycle for each issued
// certificate. A membber is created by providing it with the provider to issue
// certificates from, A request description to use when asking the provider for
// a new certificate and a handler function to run when new certificate leases
// have been issued/renewed from the provider.
//
// A new provider starts out in an initial state where it can be run from
// scratch with its Play() method. Play() is designed to be run as a goroutine
// and not doing this could result in a deadlock scenario.
type DefaultMember struct {
	context  context.Context
	cancel   func()
	doneCh   chan error
	handler  LeaseHandler
	name     string
	provider provider.Provider
	request  provider.Request
}

func (m *DefaultMember) DoneCh() <-chan error {
	return m.doneCh
}

func (m *DefaultMember) Name() string {
	return m.name
}

func NewMember(name string, provider provider.Provider, request provider.Request, handler LeaseHandler) Member {
	mem := &DefaultMember{
		doneCh:   make(chan error),
		handler:  handler,
		name:     name,
		provider: provider,
		request:  request,
	}
	mem.context, mem.cancel = context.WithCancel(context.Background())
	return mem
}

func (m *DefaultMember) Play() {
	select {
	case <-m.context.Done():
		//we return the error here because we haven't done any work and the
		//conext is alread canceled
		m.setDone(m.context.Err())
		return
	default:
	}

	log := logrus.WithFields(logrus.Fields{
		"common_name": m.request.CommonName,
	})

	log.Info("acquiring lease")
	lease, err := m.provider.Issue(&m.request)
	if err != nil {
		m.setDone(errors.Wrap(err, "issuing certificate"))
		return
	}

	log = log.WithFields(logrus.Fields{
		"lease_id":  lease.ID(),
		"lease_end": lease.Till().Format(time.RFC3339),
	})
	err = m.handler.Handle(logu.WithLogger(m.context, log), lease)
	if err != nil {
		m.setDone(errors.Wrap(err, "handling new lease"))
		return
	}

	for {
		log.Info("sleeping for next lease renewal")
		r := NewRenewer(lease)
		go r.Renew()

		select {
		case <-r.RenewCh():
			log.Info("renewing lease")
			lease, err = m.provider.Renew(lease)
			if err != nil {
				m.setDone(errors.Wrap(err, "renewing lease"))
				return
			}
			log = log.WithFields(logrus.Fields{
				"lease_id":  lease.ID(),
				"lease_end": lease.Till().Format(time.RFC3339),
			})
			log.Info("acquired new lease")
			err = m.handler.Handle(logu.WithLogger(m.context, log), lease)
			if err != nil {
				m.setDone(errors.Wrap(err, "handling renewed lease"))
				return
			}

		case err = <-r.DoneCh():
			if err != nil {
				m.setDone(errors.Wrap(err, "waiting renewer"))
				return
			}

		case <-m.context.Done():
			log.Info("stopping member lease lifecycle")
			r.Stop()
			m.setDone(nil)
			return
		}
	}
}

func (m *DefaultMember) setDone(err error) {
	m.doneCh <- err
	close(m.doneCh)
	m.cancel()
}

func (m *DefaultMember) Stop() {
	m.cancel()
}
