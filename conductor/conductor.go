package conductor

import (
	"github.com/sirupsen/logrus"
)

type Conductor struct {
	healthErr error
	members   []*MemberStatus
	memCount  int32
	watchCh   chan *MemberStatus
}

func NewConductor() *Conductor {
	return &Conductor{
		healthErr: nil,
		watchCh:   make(chan *MemberStatus),
	}
}

func (c *Conductor) AddMember(member Member) *MemberStatus {
	mstatus := NewMemberStatus(member)
	c.members = append(c.members, mstatus)
	return mstatus
}

func (c *Conductor) Play() *Conductor {
	for _, mstatus := range c.members {
		if running, err := mstatus.State(); running || err != nil {
			continue
		}
		mstatus.setState(true, nil)
		go func() {
			log := logrus.WithFields(logrus.Fields{
				"member": mstatus.member.Name(),
			})
			log.Info("playing member")
			go mstatus.member.Play()

		Outer:
			for {
				select {
				case err := <-mstatus.member.DoneCh():
					log.Info("member stopped")
					mstatus.setState(false, err)
					c.watchCh <- mstatus
					break Outer
				}
			}
		}()
	}
	return c
}

func (c *Conductor) Watch() {
	for {
		select {
		case mstatus := <-c.watchCh:
			if _, err := mstatus.State(); err != nil {
				log := logrus.WithFields(logrus.Fields{
					"member": mstatus.member.Name(),
				})
				log.Errorf("member failed: %v", err)
			}
		}
	}
}
