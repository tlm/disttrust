package conductor

import (
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	"github.com/tlmiller/disttrust/provider"
)

type Member struct {
	Provider provider.Provider
	Request  provider.Request
}

func (m *Member) Do() error {
	log := log.WithFields(log.Fields{
		"common_name": m.Request.CommonName,
	})

	log.Info("issuing certificate")
	_, err := m.Provider.Issue(&m.Request)
	if err != nil {
		return errors.Wrap(err, "issuing certificate")
	}
	return nil
}
