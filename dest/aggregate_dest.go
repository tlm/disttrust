package dest

import (
	"github.com/tlmiller/disttrust/provider"
)

type Aggregate struct {
	Dests []Dest
}

func NewAggregate(dests ...Dest) *Aggregate {
	return &Aggregate{
		Dests: dests,
	}
}

func (a *Aggregate) Send(res *provider.Response) error {
	for _, dest := range a.Dests {
		err := dest.Send(res)
		if err != nil {
			return err
		}
	}
	return nil
}
