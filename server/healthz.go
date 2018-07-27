package server

import (
	"github.com/tlmiller/disttrust/server/healthz"
)

func (a *ApiServer) AddHealthzChecks(checks ...healthz.Checker) {
	a.healthzChecks = append(a.healthzChecks, checks...)
}
